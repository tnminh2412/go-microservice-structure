package handler

import (
	"context"
	service_error "go-pkg/helper/error_helper"
	"go-pkg/service"
	"proto-repo/generated/cmm/service/cmm/core/user/core_user"
	"proto-repo/generated/cmm/service/cmm/integration/user/integration_user"
	"proto-repo/generated/pkg/nullable"
	"time"

	"github.com/tnminh2412/go-microservice-structure/service/app/integration/pkg/feature/person/user/system_user"
	"github.com/tnminh2412/go-microservice-structure/service/app/integration/user/internal/pkg/error_helper"
)

func (s *Service) Login(ctx context.Context, req *integration_user.LoginRequest) (*integration_user.LoginReply, error) {
	p, err := extractLoginParams(req)
	if err != nil {
		return nil, err
	}
	now := time.Now()

	user, err := s.userSrv.GetUser(ctx, &core_user.GetUserRequest{
		ClientID: p.ClientID,
		ID:       &nullable.NullString{Value: p.Username},
	})
	if err != nil {
		if _, isNotFoundErr := service_error.Is(err, service_error.NotFound); isNotFoundErr == true {
			return nil, error_helper.New(integration_user.ErrorMessage_InvalidUsernameOrPassword.String(), "", nil)
		}
		return nil, err
	}

	sysUser := system_user.New(ctx, p.ClientID, p.Username)
	res, err := sysUser.Login(&system_user.LoginRequest{
		Password:         p.Password,
		Otp:              p.Otp,
		LastLoginCountry: p.CountryCode,
		LastLoginIP:      p.IP,
		LastLoginTime:    now,
		UpdatedBy:        service.Name,
		UpdatedTime:      now,
	})
	if err != nil {
		return nil, err
	}
	return &integration_user.LoginReply{
		Token:                   res.Token,
		Username:                p.Username,
		IsOtpEnable:             res.IsOtpEnable,
		IsOtpSkip:               res.IsOtpSkip,
		IsRequireChangePassword: user.IsRequireChangePassword,
	}, nil

}

type loginParams struct {
	ClientID    string
	Username    string
	Password    string
	CountryCode string
	IP          string
	Otp         string
}

func extractLoginParams(req *integration_user.LoginRequest) (*loginParams, error) {
	if len(req.ClientID) == 0 {
		msg := "ClientID is required"
		return nil, error_helper.New(service_error.InvalidArgument, msg, nil)
	}
	if len(req.Username) == 0 {
		msg := "Username is required"
		return nil, error_helper.New(service_error.InvalidArgument, msg, nil)
	}
	if len(req.Password) == 0 {
		msg := "Password is required"
		return nil, error_helper.New(service_error.InvalidArgument, msg, nil)
	}

	p := &loginParams{
		ClientID:    req.ClientID,
		Username:    req.Username,
		Password:    req.Password,
		CountryCode: req.CountryCode,
		IP:          req.IP,
		Otp:         req.Otp,
	}
	return p, nil
}
