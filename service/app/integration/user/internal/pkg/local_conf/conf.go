package local_conf

import (
	"context"
	"encoding/json"
	confFile "go-pkg/conf"
	service_error "go-pkg/helper/error_helper"

	"github.com/tnminh2412/go-microservice-structure/service/app/integration/user/internal/pkg/error_helper"
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

type Locale struct {
	Language    string
	Description string
}
type PasswordRule struct {
	Regex   string
	Locales []Locale
}

type UsernameRule struct {
	Regex   string
	Locales []Locale
}

type Otp struct {
	Issuer                  string `bson:"issuer"`
	ResetTokenTtlMinutes    int64  `bson:"resetTokenTtlMinutes"`
	ResetTokenResendSeconds int64  `bson:"resetTokenResendSeconds"`
	ResetTokenParam         string `bson:"resetTokenParam"`
}

type Value struct {
	TokenTTLInHour      int32          `bson:"tokenTTLInHour"`
	PasswordRules       []PasswordRule `bson:"passwordRules"`
	UsernameRules       []UsernameRule `bson:"usernameRules"`
	MaxFailLoginAttempt int32          `bson:"maxFailLoginAttempt"`
	UsernameBlackList   []string       `bson:"usernameBlackList"`
	Otp                 *Otp           `bson:"otp"`
	Chat                *Chat          `bson:"chat"`
}

func (conf *Service) Module() string {
	appModule, _ := confFile.GetString("App", "AppModule", module)
	return appModule
}

const module = "integration_user"

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
