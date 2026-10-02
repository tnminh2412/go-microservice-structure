package handler

import (
	"context"
	service_error "go-pkg/helper/error_helper"
	"go-pkg/mongo/mongo_error"
	"proto-repo/generated/cmm/service/cmm/core/user/core_user"

	"github.com/tnminh2412/go-microservice-structure/service/app/core/user/internal/db"
	"github.com/tnminh2412/go-microservice-structure/service/app/core/user/internal/pkg/error_helper"
)

func (s *Service) CreateRole(ctx context.Context, req *core_user.CreateRoleRequest) (*core_user.CreateRoleReply, error) {
	mongo, err := s.SharedService.GetDatabase(ctx, req.ClientID)
	if err != nil {
		return nil, error_helper.New(service_error.Internal, "Server error", err)
	}

	err = db.CreateRole(ctx, mongo, &db.Role{
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
		CreatedBy:   req.CreatedBy,
		CreatedTime: req.CreatedTime.AsTime(),
		UpdatedBy:   req.CreatedBy,
		UpdatedTime: req.CreatedTime.AsTime(),
	})

	if err != nil {
		if mongo_error.IsDup(err) {
			return nil, error_helper.New(service_error.Duplicate, "Duplicated", err)
		}
		return nil, error_helper.New(service_error.Internal, "Server error", err)
	}
	return &core_user.CreateRoleReply{}, nil
}
