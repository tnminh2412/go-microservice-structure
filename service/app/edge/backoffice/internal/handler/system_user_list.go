package handler

import (
	"context"
	"go-pkg/helper/context_helper"
	"proto-repo/generated/cmm/service/cmm/edge/backoffice/edge_backoffice"
	"proto-repo/generated/cmm/service/cmm/integration/user/integration_user"

	"github.com/tnminh2412/go-microservice-structure/service/app/edge/backoffice/internal/adapter"
	"github.com/tnminh2412/go-microservice-structure/service/app/edge/backoffice/internal/handler/defined/tag"
	"github.com/tnminh2412/go-microservice-structure/service/app/edge/backoffice/internal/pkg/request_context"
)

func (s *Service) ListSystemUser(ctx context.Context, req *edge_backoffice.SystemUserListRequest) (*edge_backoffice.SystemUserListReply, error) {
	context_helper.AddTag(ctx, tag.Search, req.Search)
	context_helper.AddTag(ctx, tag.ClientID, request_context.ClientID(ctx))

	res, err := s.integrationUserServiceClient.ListSystemUser(ctx, &integration_user.SystemUserListRequest{
		ClientID: request_context.ClientID(ctx),
		Search:   req.Search,
		Statuses: req.Statuses,
		Skip:     s.getSkip(req.PageIndex, req.Limit),
		Limit:    req.Limit,
	})
	if err != nil {
		return nil, err
	}
	rep := &edge_backoffice.SystemUserListReply{
		Users: adapter.ToSystemUsers(res.Users),
	}

	return rep, nil
}
