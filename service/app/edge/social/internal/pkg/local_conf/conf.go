package local_conf

import (
	"context"
	confFile "go-pkg/conf"
	service_error "go-pkg/helper/error_helper"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/tnminh2412/go-microservice-structure/service/app/edge/social/internal/pkg/error_helper"
	"github.com/tnminh2412/go-microservice-structure/service/app/pkg/shared_conf"
)

type IService interface {
	GetConfigValue(ctx context.Context, clientID string) (*Value, error)
	Module() string
}

type Service struct {
	Service shared_conf.IService
}

func New(sharedConf shared_conf.IService) Service {
	return Service{Service: sharedConf}
}

type Serialize struct {
}

type Value struct {
}

func (conf Service) Module() string {
	appModule, _ := confFile.GetString("App", "AppModule", module)
	return appModule
}

const module = "edge_social"

func (conf Service) GetConfigValue(ctx context.Context, clientID string) (*Value, error) {
	interfaceVal, err := conf.Service.GetConf(ctx, clientID, conf.Module(), &Serialize{})
	if err != nil {
		return nil, error_helper.New(service_error.Internal, "", err)
	}

	val, ok := interfaceVal.(*Value)
	if !ok {
		return nil, error_helper.New(service_error.Internal, "", err)
	}
	return val, nil
}

func (s *Serialize) Decode(_ context.Context, result *mongo.SingleResult) (interface{}, error) {
	var _conf struct {
		Value *Value
	}
	err := result.Decode(&_conf)
	if err != nil {
		return nil, err
	}

	return _conf.Value, err
}
