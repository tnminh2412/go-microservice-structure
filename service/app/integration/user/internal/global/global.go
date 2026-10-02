package global

import (
	"context"
	"fmt"
	"go-pkg/conf"
	"go-pkg/grpc"
	"go-pkg/log"
	"go-pkg/redis"
	"proto-repo/generated/cmm/service/cmm/core/communication/core_communication"
	"proto-repo/generated/cmm/service/cmm/core/user/core_user"
	"time"

	"github.com/tnminh2412/go-microservice-structure/service/app/integration/pkg/feature/person/user/system_user"
	"github.com/tnminh2412/go-microservice-structure/service/app/integration/user/internal/adapter"
	"github.com/tnminh2412/go-microservice-structure/service/app/integration/user/internal/config"
	"github.com/tnminh2412/go-microservice-structure/service/app/integration/user/internal/config/service_config"
	"github.com/tnminh2412/go-microservice-structure/service/app/integration/user/internal/pkg/local_conf"
	"github.com/tnminh2412/go-microservice-structure/service/pkg/validator"
)

var (
	userSrv          core_user.UserServiceClient
	communicationSrv core_communication.CommunicationServiceClient
	tokenClient      redis.Client
)

func init() {
	tokenClient = redis.GetUniversal(config.TokenServer)

	userHost, _ := conf.GetString(service_config.Group, service_config.UserHostFlag, "")
	userConn, err := grpc.NewClient(userHost)
	if err != nil {
		log.Error(func() string { return fmt.Sprintf("err: %s", err.Error()) })
		panic(err)
	}
	userSrv = core_user.NewUserServiceClient(userConn)

	communicationHost, _ := conf.GetString(service_config.Group, service_config.CommunicationFlag, "")
	communicationConn, err := grpc.NewClient(communicationHost)
	if err != nil {
		log.Error(func() string { return fmt.Sprintf("err: %s", err.Error()) })
		panic(err)
	}
	communicationSrv = core_communication.NewCommunicationServiceClient(communicationConn)
}

func CoreUserClient() core_user.UserServiceClient {
	return userSrv
}

func CoreCommunication() core_communication.CommunicationServiceClient {
	return communicationSrv
}

func TokenClient() redis.Client {
	return tokenClient
}

func FnGetUsernameBlackList(localService local_conf.IService) func(ctx context.Context, clientID string) ([]string, error) {
	return func(ctx context.Context, clientID string) ([]string, error) {
		cfg, err := localService.GetConfigValue(ctx, clientID)
		if err != nil {
			return nil, err
		}
		return cfg.UsernameBlackList, nil
	}
}

func FnGetUsernameRules(localService local_conf.IService) func(ctx context.Context, clientID string) ([]validator.UsernameRule, error) {
	return func(ctx context.Context, clientID string) ([]validator.UsernameRule, error) {
		cfg, err := localService.GetConfigValue(ctx, clientID)
		if err != nil {
			return nil, err
		}

		return adapter.ToUsernameRules(cfg.UsernameRules), nil
	}
}

func FnGetPasswordRules(localService local_conf.IService) func(ctx context.Context, clientID string) ([]validator.PasswordRule, error) {
	return func(ctx context.Context, clientID string) ([]validator.PasswordRule, error) {
		cfg, err := localService.GetConfigValue(ctx, clientID)
		if err != nil {
			return nil, err
		}
		return adapter.ToPasswordRules(cfg.PasswordRules), nil
	}
}

func FnGetMaxFailLoginAttempt(localService local_conf.IService) func(ctx context.Context, clientID string) (int32, error) {
	return func(ctx context.Context, clientID string) (int32, error) {
		cfg, err := localService.GetConfigValue(ctx, clientID)
		if err != nil {
			return 0, err
		}
		return cfg.MaxFailLoginAttempt, nil
	}
}

func FnGetOtpIssuer(localService local_conf.IService) func(ctx context.Context, clientID string) (string, error) {
	return func(ctx context.Context, clientID string) (string, error) {
		cfg, err := localService.GetConfigValue(ctx, clientID)
		if err != nil {
			return "", err
		}
		return cfg.Otp.Issuer, nil
	}
}

func FnGetDefaultChatConfigSystemUser(localService local_conf.IService) func(ctx context.Context, clientID string) (*system_user.Chat, error) {
	return func(ctx context.Context, clientID string) (*system_user.Chat, error) {
		cfg, err := localService.GetConfigValue(ctx, clientID)
		if err != nil {
			return nil, err
		}
		chatConf := adapter.ToDefaultChatConfigSystemUser(cfg.Chat)
		return chatConf, nil
	}
}

func FnGetUserTokenTTL(localService local_conf.IService) func(ctx context.Context, clientID string) (time.Duration, error) {
	return func(ctx context.Context, clientID string) (time.Duration, error) {
		cfg, err := localService.GetConfigValue(ctx, clientID)
		if err != nil {
			return 0, err
		}

		return time.Duration(cfg.TokenTTLInHour) * time.Hour, nil
	}
}
