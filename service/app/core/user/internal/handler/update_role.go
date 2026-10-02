package handler

import (
	"context"
	service_error "go-pkg/helper/error_helper"
	"proto-repo/generated/cmm/service/cmm/core/user/core_user"

	"github.com/tnminh2412/go-microservice-structure/service/app/core/user/internal/db"
	"github.com/tnminh2412/go-microservice-structure/service/app/core/user/internal/pkg/error_helper"
)

func (s *Service) UpdateRole(ctx context.Context, req *core_user.UpdateRoleRequest) (*core_user.UpdateRoleReply, error) {
	mongo, err := s.SharedService.GetDatabase(ctx, req.ClientID)
	if err != nil {
		return nil, error_helper.New(service_error.Internal, "Server error", err)
	}

	matched, err := db.UpdateRole(ctx, mongo, &db.Role{
		Name:  req.Role,
		Tasks: req.Tasks,
		Locale: func() []*db.Locale {
			var locale []*db.Locale

			for _, v := range req.Locale {
				locale = append(locale, &db.Locale{
					Language:    v.Language,
					Description: v.Description,
				})
			}

			return locale
		}(),
		UpdatedBy:   req.UpdatedBy,
		UpdatedTime: req.UpdatedTime.AsTime(),
	})

	if err != nil {
		return nil, error_helper.New(service_error.Internal, "Server error", err)
	}

	if matched == 0 {
		return nil, error_helper.New(service_error.NotFound, "Not found", err)
	}

	return &core_user.UpdateRoleReply{}, nil
}
