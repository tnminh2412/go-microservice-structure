package handler

import (
	"context"
	"go-pkg/helper/context_helper"
	service_error "go-pkg/helper/error_helper"
	"proto-repo/generated/cmm/service/cmm/edge/backoffice/edge_backoffice"
	"proto-repo/generated/cmm/service/cmm/integration/user/integration_user"

	"github.com/tnminh2412/go-microservice-structure/service/app/edge/backoffice/internal/handler/defined/tag"
	"github.com/tnminh2412/go-microservice-structure/service/app/edge/backoffice/internal/pkg/access_check"
	"github.com/tnminh2412/go-microservice-structure/service/app/edge/backoffice/internal/pkg/error_helper"
	"github.com/tnminh2412/go-microservice-structure/service/app/edge/backoffice/internal/pkg/request_context"
)

func (s *Service) SetSystemUserPassword(ctx context.Context, req *edge_backoffice.SystemUserSetPasswordRequest) (*edge_backoffice.SystemUserSetPasswordReply, error) {
	context_helper.AddTag(ctx, tag.Username, req.Username)
	context_helper.AddTag(ctx, tag.ClientID, request_context.ClientID(ctx))

	err := s.accessCheck.AccessSystemUser(ctx, &access_check.AccessSystemUser{
		ClientID:       request_context.ClientID(ctx),
		TargetUsername: req.Username,
	})
	if err != nil {
		return nil, err
	}

	_, err = s.integrationUserServiceClient.SetSystemUserPassword(ctx, &integration_user.SystemUserSetPasswordRequest{
		ClientID:    request_context.ClientID(ctx),
		Username:    req.Username,
		Password:    req.Password,
		Language:    request_context.Language(ctx),
		UpdatedBy:   s.requestBy(),
		UpdatedTime: s.requestTime(),
	})
	if err != nil {
		if _, ok := service_error.Is(err, service_error.NotFound); ok {
			return nil, error_helper.New(service_error.NotFound, "", nil)
		}
		if _, ok := service_error.Is(err, integration_user.ErrorMessage_InvalidPassword.String()); ok {
			return nil, error_helper.New(edge_backoffice.ErrorMessage_InvalidPassword.String(), "", nil)
		}
		return nil, err
	}

	return &edge_backoffice.SystemUserSetPasswordReply{}, nil
}
