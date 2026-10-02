package local_conf

import (
	"context"
	"encoding/json"
	confFile "go-pkg/conf"
	service_error "go-pkg/helper/error_helper"

	"github.com/tnminh2412/go-microservice-structure/service/app/core/user/internal/pkg/error_helper"
	"github.com/tnminh2412/go-microservice-structure/service/app/pkg/shared_conf"

	"go.mongodb.org/mongo-driver/mongo"
)

type IService interface {
	GetConfigValue(ctx context.Context, clientID string) (*Value, error)
	Module() string
}

type Service struct {
	Service shared_conf.IService
}

type Serialize struct {
}

type Value struct {
}

func (conf *Service) Module() string {
	appModule, _ := confFile.GetString("App", "AppModule", module)
	return appModule
}

const module = "core_user"

func (conf *Service) GetConfigValue(ctx context.Context, clientID string) (*Value, error) {
	interfaceVal, err := conf.Service.GetConf(ctx, clientID, conf.Module(), &Serialize{})
	if err != nil {
		return nil, error_helper.New(service_error.Internal, "Configuration", err)
	}

	val, ok := interfaceVal.(*Value)
	if !ok {
		return nil, error_helper.New(service_error.Internal, "Configuration", err)
	}
	return val, nil
}

func (s *Serialize) Decode(_ context.Context, result *mongo.SingleResult) (interface{}, error) {
	var _conf struct {
		Value *Value
	}
	err := result.Decode(&_conf)
	return _conf.Value, err
}

func (s *Serialize) Deserialize(_ context.Context, bytes []byte) (interface{}, error) {
	var _conf struct {
		Value *Value
	}
	err := json.Unmarshal(bytes, &_conf)
	return _conf.Value, err
}
