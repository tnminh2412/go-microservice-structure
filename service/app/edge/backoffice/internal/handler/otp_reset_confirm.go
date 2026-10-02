package handler

import (
	"context"
	"go-pkg/helper/context_helper"
	"proto-repo/generated/cmm/service/cmm/edge/backoffice/edge_backoffice"
	"proto-repo/generated/cmm/service/cmm/integration/user/integration_user"

	"github.com/tnminh2412/go-microservice-structure/service/app/edge/backoffice/internal/handler/defined/tag"
	"github.com/tnminh2412/go-microservice-structure/service/app/edge/backoffice/internal/pkg/request_context"
)

func (s *Service) ConfirmOtpReset(ctx context.Context, req *edge_backoffice.OtpResetConfirmRequest) (*edge_backoffice.OtpResetConfirmReply, error) {
	context_helper.AddTag(ctx, tag.ClientID, request_context.ClientID(ctx))

	_, err := s.integrationUserServiceClient.ConfirmOtpReset(ctx, &integration_user.OtpResetConfirmRequest{
		ClientID:    request_context.ClientID(ctx),
		Token:       req.Token,
		RequestBy:   s.requestBy(),
		RequestTime: s.requestTime(),
	})
	if err != nil {
		return nil, err
	}

	return &edge_backoffice.OtpResetConfirmReply{}, nil
}
