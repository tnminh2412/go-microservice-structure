package handler

import (
	"context"
	"go-pkg/conf"
	"proto-repo/generated/cmm/service/cmm/core/communication/core_communication"
	"proto-repo/generated/cmm/service/cmm/core/user/core_user"
	"proto-repo/generated/cmm/service/cmm/integration/user/integration_user"
	"regexp"

	"github.com/tnminh2412/go-microservice-structure/service/app/integration/pkg/feature/person/user"
	"github.com/tnminh2412/go-microservice-structure/service/app/integration/pkg/feature/person/user/system_user"
	"github.com/tnminh2412/go-microservice-structure/service/app/integration/pkg/feature/person/user/token_user"
	"github.com/tnminh2412/go-microservice-structure/service/app/integration/user/internal/global"
	"github.com/tnminh2412/go-microservice-structure/service/app/integration/user/internal/pkg/error_helper"
	"github.com/tnminh2412/go-microservice-structure/service/app/integration/user/internal/pkg/local_conf"
	"github.com/tnminh2412/go-microservice-structure/service/app/integration/user/internal/pkg/validator"
	"github.com/tnminh2412/go-microservice-structure/service/app/pkg/shared_conf"
	"github.com/tnminh2412/go-microservice-structure/service/pkg/encrypt"
	"github.com/tnminh2412/go-microservice-structure/service/pkg/global_conf"
)

type Service struct {
	// Returns codes.Unimplemented for business RPCs stripped from this sample.
	integration_user.UnimplementedServiceServer

	GlobalService global_conf.IService
	SharedService shared_conf.IService
	LocalService  local_conf.IService

	tokenUser        *token_user.Server
	userSrv          core_user.UserServiceClient
	communicationSrv core_communication.CommunicationServiceClient
	validator        validator.Validator
	encrypter        encrypt.Encrypter
}

func New() *Service {
	globalService := global_conf.New(context.Background())
	sharedService := shared_conf.New(context.Background(), globalService)
	localService := &local_conf.Service{
		Service: sharedService,
	}

	tokenSrv := token_user.New(global.TokenClient())

	user.Init(&user.ServiceInfo{
		TokenUserService: tokenSrv,
	})

	system_user.Init(&system_user.ServiceInfo{
		UserService:              global.CoreUserClient(),
		TokenUserService:         tokenSrv,
		FnGetUsernameRules:       global.FnGetUsernameRules(localService),
		FnGetPasswordRules:       global.FnGetPasswordRules(localService),
		FnGetUsernameBlacklist:   global.FnGetUsernameBlackList(localService),
		FnGetMaxFailLoginAttempt: global.FnGetMaxFailLoginAttempt(localService),
		FnGetOtpIssuer:           global.FnGetOtpIssuer(localService),
		FnGetChatConfig:          global.FnGetDefaultChatConfigSystemUser(localService),
		FnGetTokenTTLConfig:      global.FnGetUserTokenTTL(localService),
	})

	// Secret used to sign OTP-reset tokens; set it in [Security] OtpResetSecret.
	otpResetSecret, _ := conf.GetString("Security", "OtpResetSecret", "")

	return &Service{
		GlobalService:    globalService,
		SharedService:    sharedService,
		LocalService:     localService,
		tokenUser:        tokenSrv,
		userSrv:          global.CoreUserClient(),
		communicationSrv: global.CoreCommunication(),
		validator:        validator.New(globalService, localService),
		encrypter:        encrypt.New(otpResetSecret),
	}
}

const pattern = "^\\+[0-9]+$"

func (s *Service) validatePhoneNumberPattern(phoneNumber string) error {
	reg := regexp.MustCompile(pattern)
	match := reg.MatchString(phoneNumber)
	if !match {
		msg := "PhoneNumber format invalid"
		return error_helper.New(integration_user.ErrorMessage_InvalidPhoneNumber.String(), msg, nil)
	}

	return nil
}
