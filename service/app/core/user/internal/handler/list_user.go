package handler

import (
	"context"
	service_error "go-pkg/helper/error_helper"
	"proto-repo/generated/cmm/service/cmm/core/user/core_user"

	"github.com/tnminh2412/go-microservice-structure/service/app/core/user/internal/adapter"
	"github.com/tnminh2412/go-microservice-structure/service/app/core/user/internal/db"
	"github.com/tnminh2412/go-microservice-structure/service/app/core/user/internal/pkg/error_helper"
)

func (s *Service) ListUser(ctx context.Context, req *core_user.ListUserRequest) (*core_user.ListUserReply, error) {
	mongo, err := s.SharedService.GetDatabase(ctx, req.ClientID)
	if err != nil {
		return nil, error_helper.New(service_error.Internal, "Server error", err)
	}

	p := db.ListUserParams{
		ByLevel:         req.Level != nil,
		BySkip:          req.Skip.IsValid(),
		ByLimit:         req.Limit.IsValid(),
		ByExcludeSystem: req.ExcludeSystem.IsValid(),
		ByDirectParent:  req.DirectParent.IsValid(),
		ByStatuses:      len(req.Statuses) > 0,
		ByIDs:           len(req.IDs) > 0,
		ByRoles:         len(req.Roles) > 0,
		ByExcludeRoles:  len(req.ExcludeRoles) > 0,
		ByIsSystem:      req.IsSystem.IsValid(),
		ByExcludeHidden: req.ExcludeHidden.IsValid(),

		Skip:          req.Skip.GetValue(),
		Limit:         req.Limit.GetValue(),
		IDs:           req.IDs,
		Statuses:      req.Statuses,
		DirectParent:  req.DirectParent.GetValue(),
		ExcludeSystem: req.ExcludeSystem.GetValue(),
		Level: func() *db.LevelParams {
			if req.Level == nil {
				return &db.LevelParams{}
			}
			return &db.LevelParams{
				ByEq: req.Level.Eq.IsValid(),
				ByGt: req.Level.Gt.IsValid(),
				Eq:   req.Level.Eq.GetValue(),
				Gt:   req.Level.Gt.GetValue(),
			}
		}(),
		Roles:         req.Roles,
		ExcludeRoles:  req.ExcludeRoles,
		IsSystem:      req.IsSystem.GetValue(),
		ExcludeHidden: req.ExcludeHidden.GetValue(),
		OrderBy:       req.OrderBy,
		IsAsc:         req.IsAsc,
	}

	fields := getDBExcludeMask(req.FieldMask)
	p.Mask = &db.Mask{
		Fields:    fields,
		IsExclude: true,
	}

	users, err := db.ListUser(ctx, mongo, p)
	if err != nil {
		return nil, error_helper.New(service_error.Internal, "db.ListUser error", err)
	}

	roles, err := db.ListRole(ctx, mongo, db.ListRoleParams{})
	if err != nil {
		return nil, error_helper.New(service_error.Internal, "db.ListRole error", err)
	}

	return &core_user.ListUserReply{
		Users: func() []*core_user.GetUserReply {
			reply := make([]*core_user.GetUserReply, len(users))
			for i, v := range users {
				reply[i] = adapter.ToUserReply(v, roles)
			}

			return reply
		}(),
	}, nil
}
