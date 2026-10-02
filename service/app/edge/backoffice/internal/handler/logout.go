package handler

import (
	"context"
	"go-pkg/helper/context_helper"
	"go-pkg/log/base"
	"proto-repo/generated/cmm/service/cmm/edge/backoffice/edge_backoffice"
	"proto-repo/generated/cmm/service/cmm/integration/user/integration_user"

	"github.com/tnminh2412/go-microservice-structure/service/app/edge/backoffice/internal/pkg/activity"
	"github.com/tnminh2412/go-microservice-structure/service/app/edge/backoffice/internal/pkg/request_context"
)

func (s *Service) Logout(ctx context.Context, in *edge_backoffice.LogoutRequest) (*edge_backoffice.LogoutReply, error) {
	clientID := request_context.ClientID(ctx)
	tokenInfo, err := s.userAuthorize.Get(ctx)
	if err != nil {
		return nil, err
	}
	currentUsername := tokenInfo.Username
	_, err = s.integrationUserServiceClient.Logout(ctx, &integration_user.LogoutRequest{
		ClientID: clientID,
		Token:    tokenInfo.Token,
	})
	if err != nil {
		return nil, err
	}

	err = s.logActivity(ctx, activityInfo{
		Username:   currentUsername,
		ActivityID: activity.ID.Logout,
		Data:       map[string]string{},
	})
	if err != nil {
		context_helper.Logger(ctx).Error(func() string {
			return "log activity fail"
		}, base.Error(err))
	}

	rep := &edge_backoffice.LogoutReply{}
	return rep, nil
}
