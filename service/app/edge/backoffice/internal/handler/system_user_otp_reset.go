package handler

import (
	"context"
	"go-pkg/helper/context_helper"
	"proto-repo/generated/cmm/service/cmm/edge/backoffice/edge_backoffice"
	"proto-repo/generated/cmm/service/cmm/integration/user/integration_user"

	"github.com/tnminh2412/go-microservice-structure/service/app/edge/backoffice/internal/handler/defined/tag"
	"github.com/tnminh2412/go-microservice-structure/service/app/edge/backoffice/internal/pkg/access_check"
	"github.com/tnminh2412/go-microservice-structure/service/app/edge/backoffice/internal/pkg/request_context"
)

func (s *Service) ResetSystemUserOtp(ctx context.Context, req *edge_backoffice.SystemUserResetOtpRequest) (*edge_backoffice.SystemUserResetOtpReply, error) {
	context_helper.AddTag(ctx, tag.Username, req.Username)
	context_helper.AddTag(ctx, tag.ClientID, request_context.ClientID(ctx))

	err := s.accessCheck.AccessSystemUser(ctx, &access_check.AccessSystemUser{
		ClientID:       request_context.ClientID(ctx),
		TargetUsername: req.Username,
	})
	if err != nil {
		return nil, err
	}

	_, err = s.integrationUserServiceClient.ResetSystemUserOtp(ctx, &integration_user.SystemUserResetOtpRequest{
		ClientID:    request_context.ClientID(ctx),
		Username:    req.Username,
		UpdatedBy:   s.requestBy(),
		UpdatedTime: s.requestTime(),
	})
	if err != nil {
		return nil, err
	}

	return &edge_backoffice.SystemUserResetOtpReply{}, nil
}
