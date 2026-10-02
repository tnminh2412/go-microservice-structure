package handler

import (
	"context"
	service_error "go-pkg/helper/error_helper"
	"go-pkg/mongo/mongo_error"
	"proto-repo/generated/cmm/service/cmm/core/user/core_user"

	"github.com/tnminh2412/go-microservice-structure/service/app/core/user/internal/adapter"
	"github.com/tnminh2412/go-microservice-structure/service/app/core/user/internal/db"
	"github.com/tnminh2412/go-microservice-structure/service/app/core/user/internal/pkg/error_helper"
)

func (s *Service) GetRole(ctx context.Context, req *core_user.GetRoleRequest) (*core_user.GetRoleReply, error) {
	mongo, err := s.SharedService.GetDatabase(ctx, req.ClientID)
	if err != nil {
		return nil, error_helper.New(service_error.Internal, "Server error", err)
	}

	role, err := db.GetRole(ctx, mongo, req.Role)
	if err != nil {
		if mongo_error.IsNotFound(err) {
			return nil, error_helper.New(service_error.NotFound, "Not found", err)
		}

		return nil, error_helper.New(service_error.Internal, "Server error", err)
	}

	return adapter.ToRoleReply(role), nil
}
