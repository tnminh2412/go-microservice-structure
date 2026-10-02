package shared_conf

import (
	"context"
	"encoding/json"
	"fmt"
	"go-pkg/grpc"
	"go-pkg/grpc/grpc_config"
	"go-pkg/helper/context_helper"
	"go-pkg/log/base"
	m "go-pkg/mongo"
	"go-pkg/mongo/mongo_error"
	"net/http"
	"sync"

	"github.com/tnminh2412/go-microservice-structure/service/pkg/global_conf"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const confCollectionName = "conf"

type IService interface {
	GetDatabase(ctx context.Context, clientID string) (*mongo.Database, error)
	GetConf(ctx context.Context, clientID, key string, s SerializeInterface) (interface{}, error)
	GetSharedConf(ctx context.Context, clientID string) (*SharedValue, error)
	AddRoute(server grpc.Server, module string, serializeInterface SerializeInterface)
	UpdateContact(ctx context.Context, clientID string, contactInfo *ContactInfo) error
	Shutdown() error
}

type Service struct {
	globalService   global_conf.IService
	context         context.Context
	stopFunc        context.CancelFunc
	subscribers     map[string]struct{}
	subscribersLock sync.RWMutex
	configs         map[string]map[string]interface{}
	configLock      sync.RWMutex
	connections     map[string]struct{}
}

func New(ctx context.Context, globalService global_conf.IService) IService {
	cancelContext, cancel := context.WithCancel(ctx)
	s := &Service{
		globalService: globalService,
		subscribers:   make(map[string]struct{}),
		configs:       make(map[string]map[string]interface{}),
		context:       cancelContext,
		stopFunc:      cancel,
		connections:   make(map[string]struct{}),
	}

	s.Init(ctx)

	return s
}

func (service *Service) Init(ctx context.Context) {
	clients, err := service.globalService.GetAllClientConf(ctx)
	if err != nil {
		panic(err)
	}

	opts := make([]*m.Option, 0, len(clients))
	for clientID, cfg := range clients {
		opts = append(opts, &m.Option{
			ServerName:          clientID,
			Connection:          cfg.Db.Connection,
			EncryptedConnection: cfg.Db.EncryptedConnection,
		})
		service.connections[clientID] = struct{}{}
	}
	err = m.ConnectServers(ctx, opts...)
	if err != nil {
		panic(err)
	}
}

func (service *Service) InitNewClientConnection(ctx context.Context, clientID string) error {
	cfg, err := service.globalService.GetClientConf(ctx, clientID)
	if err != nil {
		return err
	}

	opt := &m.Option{
		ServerName:          clientID,
		Connection:          cfg.Db.Connection,
		EncryptedConnection: cfg.Db.EncryptedConnection,
	}

	err = m.ConnectServers(ctx, opt)
	if err != nil {
		return err
	}

	service.connections[clientID] = struct{}{}
	return nil
}

func (service *Service) Shutdown() (err error) {
	service.subscribersLock.Lock()
	defer service.subscribersLock.Unlock()

	service.stopFunc()
	return
}

func (service *Service) GetDatabase(ctx context.Context, clientID string) (*mongo.Database, error) {
	if service.hasConnection(clientID) == false {
		err := service.InitNewClientConnection(ctx, clientID)
		if err != nil {
			return nil, err
		}
	}

	mongodb, err := service.globalService.Client(ctx, clientID)
	if err != nil {
		return nil, err
	}

	return mongodb.Database(databaseName(clientID)), nil
}

func (service *Service) GetConf(ctx context.Context, clientID, configKey string, s SerializeInterface) (interface{}, error) {
	err := service.subscribe(clientID, configKey)
	if err != nil {
		context_helper.Logger(ctx).Error(func() string {
			return fmt.Sprintf("subscribe %v %v error %v", clientID, configKey, err)
		})
	}

	//don't lock here so that warm up requests not waiting for each others
	if localConf, ok := service.getConfig(clientID, configKey); ok {
		return localConf, nil
	}

	db, err := service.GetDatabase(ctx, clientID)
	if err != nil {
		return nil, err
	}

	row := db.Collection(confCollectionName).FindOne(ctx, bson.D{{"_id", configKey}})
	result, err := s.Decode(ctx, row)

	if err != nil {
		return nil, err
	}

	service.assignConfig(result, clientID, configKey)
	return result, nil
}

func (service *Service) subscribe(clientID string, confKey string) error {
	subKey := subscribeKey(clientID, confKey)
	if check := service.check(subKey); check {
		return nil
	}

	service.subscribersLock.Lock()
	defer service.subscribersLock.Unlock()

	if _, ok := service.subscribers[subKey]; ok {
		return nil
	}

	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.M{"documentKey._id": confKey}}},
	}

	callbackFn := func(ctx context.Context, data *m.WatchEvent) {
		defer func() {
			if err := recover(); err != nil {
				context_helper.Logger(service.context).Error(func() string {
					return fmt.Sprintf("unhandled error %v %v", clientID, err)
				})
			}
		}()
		service.receiveChange(data, clientID)
	}

	err := m.Watch(service.context, clientID, databaseName(clientID), confCollectionName, pipeline, callbackFn, m.WithKey(subKey))
	if err != nil {
		return err
	}

	service.subscribers[subKey] = struct{}{}
	return nil
}

func (service *Service) receiveChange(data *m.WatchEvent, clientID string) {
	switch data.OperationType {
	case "insert":
		context_helper.Logger(service.context).Info(func() string {
			return fmt.Sprintf("insert...do nothing")
		})
	case "delete":
		context_helper.Logger(service.context).Info(func() string {
			return fmt.Sprintf("delete...%v", data.DocumentKey)
		})
		service.removeConfig(clientID, data.DocumentKey)
	case "update":
		context_helper.Logger(service.context).Info(func() string {
			return fmt.Sprintf("update...%v", data.DocumentKey)
		})
		service.removeConfig(clientID, data.DocumentKey)
	default:
		context_helper.Logger(service.context).Info(func() string {
			return fmt.Sprintf("%v...do nothing", data.OperationType)
		})
	}
}

func (service *Service) check(subKey string) bool {
	service.subscribersLock.RLock()
	defer service.subscribersLock.RUnlock()
	if _, ok := service.subscribers[subKey]; ok {
		return true
	}
	return false

}

func (service *Service) hasConnection(clientID string) bool {
	_, existed := service.connections[clientID]
	return existed
}

type SerializeInterface interface {
	Decode(ctx context.Context, result *mongo.SingleResult) (interface{}, error)
}

func (service *Service) getConfig(clientID string, configKey string) (interface{}, bool) {
	service.configLock.RLock()
	defer service.configLock.RUnlock()

	if c, ok := service.configs[clientID]; ok {
		if v, ok := c[configKey]; ok {
			return v, true
		}
	}

	return nil, false
}

func (service *Service) assignConfig(value interface{}, clientID string, configKey string) {
	service.configLock.Lock()
	defer service.configLock.Unlock()

	if _, ok := service.configs[clientID]; !ok {
		service.configs[clientID] = map[string]interface{}{}
	}

	service.configs[clientID][configKey] = value
}

func (service *Service) removeConfig(clientID string, configKey string) {
	service.configLock.Lock()
	defer service.configLock.Unlock()

	if c, ok := service.configs[clientID]; ok {
		if _, ok := c[configKey]; ok {
			delete(c, configKey)
		}
	}

}

func (service *Service) UpdateContact(ctx context.Context, clientID string, contactInfo *ContactInfo) error {
	configKey := "integration_shared"
	db, err := service.GetDatabase(ctx, clientID)
	if err != nil {
		return err
	}

	result := db.Collection(confCollectionName).FindOneAndUpdate(ctx,
		bson.D{{"_id", configKey}},
		bson.M{"$set": bson.M{"value.contactInfo.useDefault": contactInfo.UseDefault, "value.contactInfo.contacts": contactInfo.Contacts}})

	if result.Err() != nil {
		context_helper.Logger(ctx).Error(func() string {
			return fmt.Sprintf("db.Collection(confCollectionName).FindOneAndUpdate error")
		}, base.Error(result.Err()))
		return result.Err()
	}

	return nil
}

func (service *Service) AddRoute(server grpc.Server, module string, serializeInterface SerializeInterface) {
	server.AddRouter("/modules/{clientID}", func(w http.ResponseWriter, r *http.Request) {
		if !grpc_config.IsPrivateIP(r.RemoteAddr) {
			w.Write(accessDenied)
			return
		}

		vars := mux.Vars(r)
		if clientID, ok := vars["clientID"]; ok {
			configs, err := service.GetConf(r.Context(), clientID, module, serializeInterface)
			if err == nil {
				if j, err := json.MarshalIndent(configs, "", "\t"); err == nil {
					w.Write(j)
				}
			}

			if mongo_error.IsNotFound(err) {
				w.Write([]byte("{}"))
			}
		}

		return
	})
}

var accessDenied = []byte("Access Denied")

func confKey(clientID, key string) string {
	return fmt.Sprintf("_%v:conf:%v", clientID, key)
}

func databaseName(clientID string) string {
	return fmt.Sprintf("app_client_%v", clientID)
}

func subscribeKey(clientID string, key string) string {
	return fmt.Sprintf("%s_%s", clientID, key)
}
