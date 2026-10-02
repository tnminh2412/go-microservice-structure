package handler

import (
	"context"
	"go-pkg/helper/context_helper"
	service_error "go-pkg/helper/error_helper"
	"go-pkg/log/base"
	"proto-repo/generated/cmm/service/cmm/edge/backoffice/edge_backoffice"
	"proto-repo/generated/cmm/service/cmm/integration/user/integration_user"

	"github.com/tnminh2412/go-microservice-structure/service/app/edge/backoffice/internal/pkg/activity"
	"github.com/tnminh2412/go-microservice-structure/service/app/edge/backoffice/internal/pkg/error_helper"
	"github.com/tnminh2412/go-microservice-structure/service/app/edge/backoffice/internal/pkg/request_context"
)

func (s *Service) Login(ctx context.Context, in *edge_backoffice.LoginRequest) (*edge_backoffice.LoginReply, error) {
	clientID := request_context.ClientID(ctx)

	localConfig, err := s.LocalService.GetConfigValue(ctx, clientID)
	if err != nil {
		return nil, err
	}
	p := extractLoginParams(in)
	p.Username = extractUsername(localConfig.UserUsernameIgnoreCase, p.Username)

	res, err := s.integrationUserServiceClient.Login(ctx, &integration_user.LoginRequest{
		ClientID:    clientID,
		Username:    p.Username,
		Password:    p.Password,
		Otp:         p.Otp,
		CountryCode: request_context.Country(ctx),
		IP:          request_context.IP(ctx),
	})
	if err != nil {
		if _, isInvalidErr := service_error.Is(err, integration_user.ErrorMessage_InvalidUsernameOrPassword.String()); isInvalidErr == true {
			return nil, error_helper.New(edge_backoffice.ErrorMessage_InvalidUsernameOrPassword.String(), "", nil)
		}
		if _, ok := service_error.Is(err, integration_user.ErrorMessage_MaxLoginAttemptExceed.String()); ok {
			return nil, error_helper.New(edge_backoffice.ErrorMessage_MaxLoginAttemptExceed.String(), "", nil)
		}
		if s, ok := service_error.Is(err, integration_user.ErrorMessage_InvalidStatusLocked.String()); ok {
			return nil, error_helper.New(edge_backoffice.ErrorMessage_InvalidStatusLocked.String(), s.Message, nil)
		}
		if s, ok := service_error.Is(err, integration_user.ErrorMessage_InvalidStatusSuspended.String()); ok {
			return nil, error_helper.New(edge_backoffice.ErrorMessage_InvalidStatusSuspended.String(), s.Message, nil)
		}
		if s, ok := service_error.Is(err, integration_user.ErrorMessage_InvalidStatus.String()); ok {
			return nil, error_helper.New(edge_backoffice.ErrorMessage_InvalidStatus.String(), s.Message, nil)
		}
		return nil, err
	}

	err = s.logActivity(ctx, activityInfo{
		Username:   res.Username,
		ActivityID: activity.ID.Login,
		Data:       map[string]string{},
	})
	if err != nil {
		context_helper.Logger(ctx).Error(func() string {
			return "log activity fail"
		}, base.Error(err))
	}

	rep := &edge_backoffice.LoginReply{
		Token:                   res.Token,
		IsOtpEnable:             res.IsOtpEnable,
		SkipOtp:                 res.IsOtpSkip,
		IsRequireChangePassword: res.IsRequireChangePassword,
	}
	return rep, nil
}

type loginParams struct {
	Username string
	Password string
	Otp      string
}

func extractLoginParams(req *edge_backoffice.LoginRequest) *loginParams {
	return &loginParams{
		Username: req.Username,
		Password: req.Password,
		Otp:      req.Otp,
	}
}
