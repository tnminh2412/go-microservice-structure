package handler

import (
	"context"
	"go-pkg/helper/context_helper"
	service_error "go-pkg/helper/error_helper"
	"proto-repo/generated/cmm/service/cmm/integration/user/integration_user"

	"github.com/tnminh2412/go-microservice-structure/service/app/integration/pkg/feature/person/user/system_user"
	"github.com/tnminh2412/go-microservice-structure/service/app/integration/user/internal/handler/defined/tag"
	"github.com/tnminh2412/go-microservice-structure/service/app/integration/user/internal/pkg/error_helper"

	"github.com/golang/protobuf/ptypes/timestamp"
)

func (s *Service) SetSystemUserPassword(ctx context.Context, req *integration_user.SystemUserSetPasswordRequest) (*integration_user.SystemUserSetPasswordReply, error) {
	context_helper.AddTag(ctx, tag.ClientID, req.ClientID)
	context_helper.AddTag(ctx, tag.Username, req.Username)

	p, err := s.extractSetSystemUserPasswordParams(req)
	if err != nil {
		return nil, err
	}

	sysUser := system_user.New(ctx, p.ClientID, p.Username)
	err = sysUser.SetPassword(&system_user.SetPasswordRequest{
		Password:    p.Password,
		Language:    p.Language,
		UpdatedBy:   p.UpdatedBy,
		UpdatedTime: p.UpdatedTime.AsTime(),
	})
	if err != nil {
		return nil, err
	}
	res := &integration_user.SystemUserSetPasswordReply{}
	return res, nil
}

type setSystemUserPasswordParams struct {
	ClientID    string
	Username    string
	Password    string
	Language    string
	UpdatedBy   string
	UpdatedTime *timestamp.Timestamp
}

func (s *Service) extractSetSystemUserPasswordParams(req *integration_user.SystemUserSetPasswordRequest) (*setSystemUserPasswordParams, error) {
	if len(req.ClientID) == 0 {
		msg := "ClientID is required"
		return nil, error_helper.New(service_error.InvalidArgument, msg, nil)
	}
	if len(req.Username) == 0 {
		msg := "Username is required"
		return nil, error_helper.New(service_error.InvalidArgument, msg, nil)
	}

	p := &setSystemUserPasswordParams{
		ClientID:    req.ClientID,
		Username:    req.Username,
		Password:    req.Password,
		Language:    req.Language,
		UpdatedBy:   req.UpdatedBy,
		UpdatedTime: req.UpdatedTime,
	}
	return p, nil
}
