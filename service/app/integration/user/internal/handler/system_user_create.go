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

func (s *Service) CreateSystemUser(ctx context.Context, req *integration_user.SystemUserCreateRequest) (*integration_user.SystemUserCreateReply, error) {
	context_helper.AddTag(ctx, tag.ClientID, req.ClientID)
	context_helper.AddTag(ctx, tag.Username, req.Username)

	p, err := s.extractCreateSystemUserParams(req)
	if err != nil {
		return nil, err
	}

	user := system_user.New(ctx, p.ClientID, p.Username)
	err = user.Create(&system_user.CreateRequest{
		Password:    p.Password,
		FirstName:   p.FirstName,
		LastName:    p.LastName,
		DisplayName: p.DisplayName,
		CreatedBy:   p.CreatedBy,
		CreatedTime: p.CreatedTime.AsTime(),
		Roles:       p.Roles,
		Language:    p.Language,
		AccessLines: p.AccessLines,
	})
	if err != nil {
		return nil, err
	}
	return &integration_user.SystemUserCreateReply{}, nil
}

type systemUserCreateParams struct {
	ClientID    string
	Username    string
	DisplayName string
	Password    string
	FirstName   string
	LastName    string
	Roles       []string
	Language    string
	AccessLines []string
	CreatedBy   string
	CreatedTime *timestamp.Timestamp
}

func (s *Service) extractCreateSystemUserParams(req *integration_user.SystemUserCreateRequest) (*systemUserCreateParams, error) {
	if len(req.ClientID) == 0 {
		msg := "ClientID is required"
		return nil, error_helper.New(service_error.InvalidArgument, msg, nil)
	}
	if len(req.Username) == 0 {
		msg := "Username is required"
		return nil, error_helper.New(integration_user.ErrorMessage_InvalidUsername.String(), msg, nil)
	}
	if len(req.Roles) == 0 {
		msg := "Roles is required"
		return nil, error_helper.New(integration_user.ErrorMessage_InvalidRoleIsRequired.String(), msg, nil)
	}

	p := &systemUserCreateParams{
		ClientID:    req.ClientID,
		Username:    req.Username,
		DisplayName: req.DisplayName,
		Password:    req.Password,
		Roles:       req.Roles,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Language:    req.Language,
		AccessLines: req.AccessLines,
		CreatedBy:   req.CreatedBy,
		CreatedTime: req.CreatedTime,
	}
	return p, nil
}
