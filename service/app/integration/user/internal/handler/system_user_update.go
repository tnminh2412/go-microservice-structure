package handler

import (
	"context"
	"go-pkg/helper/context_helper"
	service_error "go-pkg/helper/error_helper"
	"proto-repo/generated/cmm/service/cmm/integration/user/integration_user"
	"proto-repo/generated/pkg/nullable"

	"github.com/tnminh2412/go-microservice-structure/service/app/integration/pkg/feature/person/user/system_user"
	"github.com/tnminh2412/go-microservice-structure/service/app/integration/user/internal/handler/defined/tag"
	"github.com/tnminh2412/go-microservice-structure/service/app/integration/user/internal/pkg/error_helper"

	"github.com/golang/protobuf/ptypes/timestamp"
)

func (s *Service) UpdateSystemUser(ctx context.Context, req *integration_user.SystemUserUpdateRequest) (*integration_user.SystemUserUpdateReply, error) {
	context_helper.AddTag(ctx, tag.ClientID, req.ClientID)
	context_helper.AddTag(ctx, tag.Username, req.Username)

	p, err := extractUpdateSystemUserParams(req)
	if err != nil {
		return nil, err
	}
	sysUser := system_user.New(ctx, p.ClientID, p.Username)
	err = sysUser.Update(&system_user.UpdateRequest{
		FirstName:   p.FirstName,
		LastName:    p.LastName,
		Roles:       p.Roles,
		DisplayName: p.DisplayName,
		UpdatedBy:   p.UpdatedBy,
		UpdatedTime: p.UpdatedTime.AsTime(),
		Email:       p.Email,
	})
	if err != nil {
		return nil, err
	}
	return &integration_user.SystemUserUpdateReply{}, nil
}

type systemUserUpdateParams struct {
	ClientID    string
	Username    string
	FirstName   *nullable.NullString
	LastName    *nullable.NullString
	DisplayName *nullable.NullString
	Roles       []string
	UpdatedBy   string
	UpdatedTime *timestamp.Timestamp
	Email       *nullable.NullString
}

func extractUpdateSystemUserParams(req *integration_user.SystemUserUpdateRequest) (*systemUserUpdateParams, error) {
	if len(req.ClientID) == 0 {
		msg := "ClientID is required"
		return nil, error_helper.New(service_error.InvalidArgument, msg, nil)
	}
	if len(req.Username) == 0 {
		msg := "Username is required"
		return nil, error_helper.New(service_error.InvalidArgument, msg, nil)
	}

	p := &systemUserUpdateParams{
		ClientID:    req.ClientID,
		Username:    req.Username,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		DisplayName: req.DisplayName,
		Roles: func() []string {
			if req.Roles == nil {
				return nil
			}
			return req.Roles.List
		}(),
		UpdatedBy:   req.UpdatedBy,
		UpdatedTime: req.UpdatedTime,
		Email:       req.Email,
	}
	return p, nil
}
