package global_conf

import (
	"context"
	"encoding/json"
	"go-pkg/grpc"
	"go-pkg/grpc/grpc_config"
	"go-pkg/helper/context_helper"
	"go-pkg/log/base"
	m "go-pkg/mongo"
	"net/http"
	"sync"

	"github.com/gorilla/mux"

	"go.mongodb.org/mongo-driver/mongo"
)

const globalCollectionName = "app_global"

type IService interface {
	GetClientConf(context.Context, string) (*Conf, error)
	GetAllClientConf(context.Context) (map[string]*Conf, error)
	Client(context.Context, string) (*mongo.Client, error)
	SSOClient(context.Context, string) (*mongo.Client, error)
	GetTemplate(context.Context, string) (*Template, error)
	GetClientID(ctx context.Context) ([]string, error)
	AddRoute(server grpc.Server)
	Shutdown() error
}

type Service struct {
	context         context.Context
	stopFunc        context.CancelFunc
	subscribers     map[string]struct{}
	subscribersLock sync.RWMutex
	configs         map[string]*Conf
	configLock      sync.RWMutex
	templates       map[string]*Template
	templateLock    sync.RWMutex
}

func New(ctx context.Context) IService {
	cfgCancelContext, cfgCancelFunc := context.WithCancel(ctx)
	service := &Service{
		subscribers: make(map[string]struct{}),
		configs:     make(map[string]*Conf),
		templates:   make(map[string]*Template),
		context:     cfgCancelContext,
		stopFunc:    cfgCancelFunc,
	}

	if err := service.subscribeClient(); err != nil {
		context_helper.Logger(service.context).Error(func() string {
			return "subscribeClient error"
		}, base.Error(err))
	}

	return service
}

func (service *Service) Client(ctx context.Context, clientID string) (*mongo.Client, error) {
	conf, err := service.GetClientConf(ctx, clientID)
	if err != nil {
		return nil, err
	}
	return conf.Client()
}

func (service *Service) GetClientID(ctx context.Context) ([]string, error) {
	globalConfigs, err := service.GetAllClientConf(ctx)
	if err != nil {
		context_helper.Logger(ctx).Error(func() string {
			return "get conf err"
		}, base.Error(err))
		return nil, err
	}

	var clientIDs []string
	for clientID := range globalConfigs {
		clientIDs = append(clientIDs, clientID)
	}

	return clientIDs, nil
}

func (service *Service) SSOClient(ctx context.Context, clientID string) (*mongo.Client, error) {
	conf, err := service.GetClientConf(ctx, clientID)
	if err != nil {
		return nil, err
	}
	return conf.SSOClient()
}

func (service *Service) Shutdown() (err error) {
	service.subscribersLock.Lock()
	defer service.subscribersLock.Unlock()

	service.stopFunc()
	return
}

func (service *Service) subscribe(collection string, callBackFun func(context.Context, *m.WatchEvent)) error {
	subKey := subscribeKey(collection)
	if service.check(subKey) {
		return nil
	}

	service.subscribersLock.Lock()
	defer service.subscribersLock.Unlock()

	if _, ok := service.subscribers[subKey]; ok == true {
		return nil
	}

	pipeline := mongo.Pipeline{}

	err := m.Watch(service.context, "", globalCollectionName, collection, pipeline, callBackFun, m.WithKey(collection))
	if err != nil {
		return err
	}
	service.subscribers[subKey] = struct{}{}
	return nil
}

func (service *Service) check(watchKey string) bool {
	service.subscribersLock.RLock()
	defer service.subscribersLock.RUnlock()
	_, existed := service.subscribers[watchKey]
	return existed
}

func subscribeKey(collection string) string {
	return collection
}

func (service *Service) AddRoute(server grpc.Server) {
	server.AddRouter("/clients", func(w http.ResponseWriter, r *http.Request) {
		if !grpc_config.IsPrivateIP(r.RemoteAddr) {
			_, _ = w.Write(accessDenied)
			return
		}
		configs, err := service.GetAllClientConf(r.Context())
		if err == nil {
			if j, err := json.MarshalIndent(configs, "", "\t"); err == nil {
				_, _ = w.Write(j)
			}
		}
		return
	})

	server.AddRouter("/clients/{clientID}", func(w http.ResponseWriter, r *http.Request) {
		if !grpc_config.IsPrivateIP(r.RemoteAddr) {
			_, _ = w.Write(accessDenied)
			return
		}
		vars := mux.Vars(r)
		if clientID, ok := vars["clientID"]; ok {

			configs, err := service.GetClientConf(r.Context(), clientID)
			if err == nil {
				if j, err := json.MarshalIndent(configs, "", "\t"); err == nil {
					_, _ = w.Write(j)
				}
			}
		}
		return
	})

}

var accessDenied = []byte("Access Denied")
