package handler

import (
	"context"
	"go-pkg/helper/context_helper"
	"proto-repo/generated/cmm/service/cmm/edge/backoffice/edge_backoffice"
	"proto-repo/generated/cmm/service/cmm/integration/user/integration_user"

	"github.com/tnminh2412/go-microservice-structure/service/app/edge/backoffice/internal/handler/defined/tag"
	"github.com/tnminh2412/go-microservice-structure/service/app/edge/backoffice/internal/pkg/request_context"
)

func (s *Service) ListRole(ctx context.Context, in *edge_backoffice.RoleListRequest) (*edge_backoffice.RoleListReply, error) {
	context_helper.AddTag(ctx, tag.ClientID, request_context.ClientID(ctx))

	res, err := s.integrationUserServiceClient.ListRole(ctx, &integration_user.ListRoleRequest{
		ClientID: request_context.ClientID(ctx),
	})
	if err != nil {
		return nil, err
	}

	roles := make([]*edge_backoffice.Role, 0)
	for _, r := range res.Roles {
		roles = append(roles, &edge_backoffice.Role{
			Name: r.Role,
			Locales: func() []*edge_backoffice.Locale {
				list := make([]*edge_backoffice.Locale, 0)
				for _, l := range r.Locales {
					list = append(list, &edge_backoffice.Locale{
						Language:    l.Language,
						Description: l.Description,
					})
				}
				return list
			}(),
		})
	}

	return &edge_backoffice.RoleListReply{
		Roles: roles,
	}, nil
}
