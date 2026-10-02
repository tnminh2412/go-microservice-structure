package handler

import (
	"context"
	"go-pkg/helper/context_helper"
	"proto-repo/generated/cmm/service/cmm/edge/backoffice/edge_backoffice"
	"proto-repo/generated/cmm/service/cmm/integration/user/integration_user"

	"github.com/tnminh2412/go-microservice-structure/service/app/edge/backoffice/internal/handler/defined/tag"
	"github.com/tnminh2412/go-microservice-structure/service/app/edge/backoffice/internal/pkg/request_context"
)

func (s *Service) RequestOtpReset(ctx context.Context, req *edge_backoffice.OtpResetRequest) (*edge_backoffice.OtpResetReply, error) {
	context_helper.AddTag(ctx, tag.Email, req.Email)
	context_helper.AddTag(ctx, tag.ClientID, request_context.ClientID(ctx))

	_, err := s.integrationUserServiceClient.RequestOtpReset(ctx, &integration_user.OtpResetRequest{
		ClientID:    request_context.ClientID(ctx),
		Email:       req.Email,
		Language:    request_context.Language(ctx),
		RequestBy:   s.requestBy(),
		RequestTime: s.requestTime(),
		RedirectUrl: req.RedirectUrl,
	})
	if err != nil {
		return nil, err
	}

	return &edge_backoffice.OtpResetReply{}, nil
}
