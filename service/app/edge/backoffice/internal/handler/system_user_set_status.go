package handler

import (
	"context"
	"go-pkg/helper/context_helper"
	service_error "go-pkg/helper/error_helper"
	"go-pkg/log/base"
	"proto-repo/generated/cmm/service/cmm/edge/backoffice/edge_backoffice"
	"proto-repo/generated/cmm/service/cmm/integration/user/integration_user"

	"github.com/tnminh2412/go-microservice-structure/service/app/edge/backoffice/internal/handler/defined/tag"
	"github.com/tnminh2412/go-microservice-structure/service/app/edge/backoffice/internal/pkg/access_check"
	"github.com/tnminh2412/go-microservice-structure/service/app/edge/backoffice/internal/pkg/activity"
	"github.com/tnminh2412/go-microservice-structure/service/app/edge/backoffice/internal/pkg/error_helper"
	"github.com/tnminh2412/go-microservice-structure/service/app/edge/backoffice/internal/pkg/request_context"
)

func (s *Service) SetSystemUserStatus(ctx context.Context, req *edge_backoffice.SystemUserSetStatusRequest) (*edge_backoffice.SystemUserSetStatusReply, error) {
	context_helper.AddTag(ctx, tag.Username, req.Username)
	context_helper.AddTag(ctx, tag.ClientID, request_context.ClientID(ctx))
	context_helper.AddTag(ctx, tag.Status, req.Status)

	err := s.accessCheck.AccessSystemUser(ctx, &access_check.AccessSystemUser{
		ClientID:       request_context.ClientID(ctx),
		TargetUsername: req.Username,
	})
	if err != nil {
		return nil, err
	}

	tokenInfo, err := s.userAuthorize.Get(ctx)
	if err != nil {
		return nil, err
	}

	_, err = s.integrationUserServiceClient.SetSystemUserStatus(ctx, &integration_user.SystemUserSetStatusRequest{
		ClientID:    request_context.ClientID(ctx),
		Username:    req.Username,
		Status:      req.Status,
		UpdatedBy:   s.requestBy(),
		UpdatedTime: s.requestTime(),
	})
	if err != nil {
		if _, ok := service_error.Is(err, service_error.NotFound); ok {
			return nil, error_helper.New(service_error.NotFound, "", nil)
		}
		return nil, err
	}

	err = s.logActivity(ctx, activityInfo{
		Username:   tokenInfo.Username,
		ActivityID: activity.ID.SetStatus,
		Data: map[string]string{
			activity.Data.SystemUser: req.Username,
			activity.Data.Status:     req.Status,
		},
	})
	if err != nil {
		context_helper.Logger(ctx).Error(func() string {
			return "log activity fail"
		}, base.Error(err))
	}

	return &edge_backoffice.SystemUserSetStatusReply{}, nil
}
