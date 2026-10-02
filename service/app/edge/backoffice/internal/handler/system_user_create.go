package handler

import (
	"context"
	"go-pkg/helper/context_helper"
	service_error "go-pkg/helper/error_helper"
	"proto-repo/generated/cmm/service/cmm/edge/backoffice/edge_backoffice"
	"proto-repo/generated/cmm/service/cmm/integration/user/integration_user"
	"strings"

	"github.com/tnminh2412/go-microservice-structure/service/app/edge/backoffice/internal/handler/defined/tag"
	"github.com/tnminh2412/go-microservice-structure/service/app/edge/backoffice/internal/pkg/error_helper"
	"github.com/tnminh2412/go-microservice-structure/service/app/edge/backoffice/internal/pkg/request_context"
)

func (s *Service) CreateSystemUser(ctx context.Context, req *edge_backoffice.SystemUserCreateRequest) (*edge_backoffice.SystemUserCreateReply, error) {
	context_helper.AddTag(ctx, tag.Username, req.Username)
	context_helper.AddTag(ctx, tag.ClientID, request_context.ClientID(ctx))
	req.Username = strings.ToLower(req.Username)

	tokenInfo, err := s.userAuthorize.Get(ctx)
	if err != nil {
		return nil, err
	}

	_, err = s.integrationUserServiceClient.CreateSystemUser(ctx, &integration_user.SystemUserCreateRequest{
		ClientID:    request_context.ClientID(ctx),
		Username:    req.Username,
		DisplayName: req.Username,
		Password:    req.Password,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Language:    request_context.Language(ctx),
		Roles:       req.Roles,
		AccessLines: tokenInfo.AccessLines,
		CreatedBy:   s.requestBy(),
		CreatedTime: s.requestTime(),
	})
	if err != nil {
		if _, ok := service_error.Is(err, integration_user.ErrorMessage_InvalidPassword.String()); ok {
			return nil, error_helper.New(edge_backoffice.ErrorMessage_InvalidPassword.String(), "", nil)
		}
		if _, ok := service_error.Is(err, integration_user.ErrorMessage_InvalidUsername.String()); ok {
			return nil, error_helper.New(edge_backoffice.ErrorMessage_InvalidUsername.String(), "", nil)
		}
		if _, ok := service_error.Is(err, integration_user.ErrorMessage_InvalidUsernameIsRequired.String()); ok {
			return nil, error_helper.New(edge_backoffice.ErrorMessage_InvalidUsernameIsRequired.String(), "", nil)
		}
		if _, ok := service_error.Is(err, integration_user.ErrorMessage_InvalidRoleIsRequired.String()); ok {
			return nil, error_helper.New(edge_backoffice.ErrorMessage_InvalidRoleIsRequired.String(), "", nil)
		}
		return nil, err
	}

	return &edge_backoffice.SystemUserCreateReply{}, nil
}
