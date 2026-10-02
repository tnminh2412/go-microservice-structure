package handler

import (
	"context"
	"go-pkg/helper/context_helper"
	service_error "go-pkg/helper/error_helper"
	"proto-repo/generated/cmm/service/cmm/core/user/core_user"
	"proto-repo/generated/cmm/service/cmm/integration/user/integration_user"

	"github.com/tnminh2412/go-microservice-structure/service/app/integration/user/internal/handler/defined/tag"
	"github.com/tnminh2412/go-microservice-structure/service/app/integration/user/internal/pkg/error_helper"
)

func (s *Service) ListRole(ctx context.Context, req *integration_user.ListRoleRequest) (*integration_user.ListRoleReply, error) {
	context_helper.AddTag(ctx, tag.ClientID, req.ClientID)

	p, err := s.extractListUserRoleParams(req)
	if err != nil {
		return nil, err
	}
	roles, err := s.listRoles(ctx, p.ClientID)
	if err != nil {
		return nil, err
	}
	rep := &integration_user.ListRoleReply{
		Roles: roles,
	}
	return rep, nil
}

type userRoleListParams struct {
	ClientID string
}

func (s *Service) extractListUserRoleParams(req *integration_user.ListRoleRequest) (*userRoleListParams, error) {
	if len(req.ClientID) == 0 {
		msg := "ClientID is required"
		return nil, error_helper.New(service_error.InvalidArgument, msg, nil)
	}

	p := &userRoleListParams{
		ClientID: req.ClientID,
	}
	return p, nil
}

func (s *Service) listRoles(ctx context.Context, clientID string) ([]*integration_user.TaskRole, error) {
	params := &core_user.ListRoleRequest{
		ClientID:      clientID,
		ExcludeSystem: true,
	}
	res, err := s.userSrv.ListRole(ctx, params)
	if err != nil {
		return nil, err
	}
	roles := make([]*integration_user.TaskRole, 0)
	for _, r := range res.Roles {
		roles = append(roles, &integration_user.TaskRole{
			Role:  r.Role,
			Tasks: r.Tasks,
			Locales: func() []*integration_user.Locale {
				list := make([]*integration_user.Locale, 0)
				for _, l := range r.Locale {
					list = append(list, &integration_user.Locale{
						Language:    l.Language,
						Description: l.Description,
					})
				}
				return list
			}(),
		})
	}
	return roles, nil
}
