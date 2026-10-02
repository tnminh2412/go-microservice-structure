package local_conf

import (
	"context"
	confFile "go-pkg/conf"
	service_error "go-pkg/helper/error_helper"
	"time"

	"github.com/tnminh2412/go-microservice-structure/service/app/edge/backoffice/internal/pkg/error_helper"
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

func New(sharedConf shared_conf.IService) Service {
	return Service{Service: sharedConf}
}

type Serialize struct {
}

// Value is the per-tenant config of this service, stored in MongoDB and
// hot-reloaded by shared_conf.
type Value struct {
	IsMaintenance          bool       `bson:"isMaintenance"`
	MaintenanceStartTime   *time.Time `bson:"maintenanceStartTime"`
	MaintenanceEndTime     *time.Time `bson:"maintenanceEndTime"`
	UserUsernameIgnoreCase bool       `bson:"userUsernameIgnoreCase"`
}

func (conf Service) Module() string {
	appModule, _ := confFile.GetString("App", "AppModule", module)
	return appModule
}

const module = "edge_backoffice"

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

	return _conf.Value, nil
}
