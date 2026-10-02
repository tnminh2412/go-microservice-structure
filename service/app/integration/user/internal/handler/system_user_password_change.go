package handler

import (
	"context"
	"go-pkg/helper/context_helper"
	service_error "go-pkg/helper/error_helper"
	"go-pkg/service"
	"proto-repo/generated/cmm/service/cmm/integration/user/integration_user"
	"time"

	"github.com/tnminh2412/go-microservice-structure/service/app/integration/pkg/feature/person/user/system_user"
	"github.com/tnminh2412/go-microservice-structure/service/app/integration/user/internal/handler/defined/tag"
	"github.com/tnminh2412/go-microservice-structure/service/app/integration/user/internal/pkg/error_helper"
)

func (s *Service) SystemUserChangePassword(ctx context.Context, req *integration_user.SystemUserChangePasswordRequest) (*integration_user.SystemUserChangePasswordReply, error) {
	context_helper.AddTag(ctx, tag.ClientID, req.ClientID)
	context_helper.AddTag(ctx, tag.Username, req.Username)

	err := s.validateSystemUserChangePasswordParams(req)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	sysUser := system_user.New(ctx, req.ClientID, req.Username)
	err = sysUser.ChangePassword(&system_user.ChangePasswordRequest{
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
		Language:    req.Language,
		UpdatedBy:   service.Name,
		UpdatedTime: now,
	})
	if err != nil {
		return nil, err
	}

	return &integration_user.SystemUserChangePasswordReply{}, nil
}

func (s *Service) validateSystemUserChangePasswordParams(req *integration_user.SystemUserChangePasswordRequest) error {
	if len(req.ClientID) == 0 {
		return error_helper.New(service_error.InvalidArgument, "ClientID is required", nil)
	}
	if len(req.Username) == 0 {
		return error_helper.New(service_error.InvalidArgument, "Username is required", nil)
	}
	if len(req.OldPassword) == 0 {
		return error_helper.New(service_error.InvalidArgument, "OldPassword is required", nil)
	}
	if len(req.NewPassword) == 0 {
		return error_helper.New(service_error.InvalidArgument, "NewPassword is required", nil)
	}
	return nil
}
