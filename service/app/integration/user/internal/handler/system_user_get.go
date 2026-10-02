package handler

import (
	"context"
	"go-pkg/helper/context_helper"
	service_error "go-pkg/helper/error_helper"
	"proto-repo/generated/cmm/service/cmm/integration/user/integration_user"

	"github.com/tnminh2412/go-microservice-structure/service/app/integration/pkg/feature/person/user/system_user"
	"github.com/tnminh2412/go-microservice-structure/service/app/integration/user/internal/adapter"
	"github.com/tnminh2412/go-microservice-structure/service/app/integration/user/internal/handler/defined/tag"
	"github.com/tnminh2412/go-microservice-structure/service/app/integration/user/internal/pkg/error_helper"
)

func (s *Service) GetSystemUser(ctx context.Context, req *integration_user.SystemUserGetRequest) (*integration_user.SystemUser, error) {
	context_helper.AddTag(ctx, tag.ClientID, req.ClientID)
	context_helper.AddTag(ctx, tag.Username, req.Username)

	p, err := extractGetSystemUserParams(req)
	if err != nil {
		return nil, err
	}

	user := system_user.New(ctx, p.ClientID, p.Username)
	sysUser, err := user.Get()
	if err != nil {
		return nil, err
	}

	return adapter.ToGetSystemUser(sysUser), nil
}

type systemUserGetParams struct {
	ClientID string
	Username string
}

func extractGetSystemUserParams(req *integration_user.SystemUserGetRequest) (*systemUserGetParams, error) {
	if len(req.ClientID) == 0 {
		msg := "ClientID is required"
		return nil, error_helper.New(service_error.InvalidArgument, msg, nil)
	}
	if len(req.Username) == 0 {
		msg := "Username is required"
		return nil, error_helper.New(service_error.InvalidArgument, msg, nil)
	}

	p := &systemUserGetParams{
		ClientID: req.ClientID,
		Username: req.Username,
	}
	return p, nil
}
