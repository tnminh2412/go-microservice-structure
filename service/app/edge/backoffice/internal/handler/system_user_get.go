package handler

import (
	"context"
	"go-pkg/helper/context_helper"
	service_error "go-pkg/helper/error_helper"
	"proto-repo/generated/cmm/service/cmm/edge/backoffice/edge_backoffice"
	"proto-repo/generated/cmm/service/cmm/integration/user/integration_user"

	"github.com/tnminh2412/go-microservice-structure/service/app/edge/backoffice/internal/adapter"
	"github.com/tnminh2412/go-microservice-structure/service/app/edge/backoffice/internal/handler/defined/tag"
	"github.com/tnminh2412/go-microservice-structure/service/app/edge/backoffice/internal/pkg/access_check"
	"github.com/tnminh2412/go-microservice-structure/service/app/edge/backoffice/internal/pkg/error_helper"
	"github.com/tnminh2412/go-microservice-structure/service/app/edge/backoffice/internal/pkg/request_context"
)

func (s *Service) GetSystemUser(ctx context.Context, req *edge_backoffice.SystemUserGetRequest) (*edge_backoffice.SystemUser, error) {
	context_helper.AddTag(ctx, tag.Username, req.Username)
	context_helper.AddTag(ctx, tag.ClientID, request_context.ClientID(ctx))

	err := s.accessCheck.AccessSystemUser(ctx, &access_check.AccessSystemUser{
		ClientID:       request_context.ClientID(ctx),
		TargetUsername: req.Username,
	})
	if err != nil {
		return nil, err
	}

	u, err := s.integrationUserServiceClient.GetSystemUser(ctx, &integration_user.SystemUserGetRequest{
		ClientID: request_context.ClientID(ctx),
		Username: req.Username,
	})
	if err != nil {
		if _, ok := service_error.Is(err, service_error.NotFound); ok {
			return nil, error_helper.New(service_error.NotFound, "user not found", nil)
		}
		return nil, err
	}

	return adapter.ToSystemUser(u), nil
}
