package handler

import (
	"context"
	"go-pkg/helper/array_helper"
	service_error "go-pkg/helper/error_helper"
	"proto-repo/generated/cmm/service/cmm/core/user/core_user"

	"github.com/tnminh2412/go-microservice-structure/service/app/core/user/internal/adapter"
	"github.com/tnminh2412/go-microservice-structure/service/app/core/user/internal/db"
	"github.com/tnminh2412/go-microservice-structure/service/app/core/user/internal/pkg/error_helper"
)

func (s *Service) ListRole(ctx context.Context, req *core_user.ListRoleRequest) (*core_user.ListRoleReply, error) {

	mongo, err := s.SharedService.GetDatabase(ctx, req.ClientID)
	if err != nil {
		return nil, error_helper.New(service_error.Internal, "Server error", err)
	}

	roles, err := db.ListRole(ctx, mongo, db.ListRoleParams{
		ExcludeSystem: req.ExcludeSystem,
	})
	if err != nil {
		return nil, error_helper.New(service_error.Internal, "Server error", err)
	}

	return &core_user.ListRoleReply{
		Roles: func() []*core_user.GetRoleReply {
			reply := make([]*core_user.GetRoleReply, 0, len(roles))
			for _, r := range roles {
				if len(req.IDs) == 0 || array_helper.Contains(req.IDs, r.Name, false) == true {
					reply = append(reply, adapter.ToRoleReply(r))
				}
			}
			return reply
		}(),
	}, nil
}
