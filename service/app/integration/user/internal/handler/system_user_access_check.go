package handler

import (
	"context"
	"go-pkg/helper/context_helper"
	service_error "go-pkg/helper/error_helper"
	"proto-repo/generated/cmm/service/cmm/integration/user/integration_user"

	"github.com/tnminh2412/go-microservice-structure/service/app/integration/pkg/feature/person/user/system_user"
	"github.com/tnminh2412/go-microservice-structure/service/app/integration/user/internal/handler/defined/tag"
	"github.com/tnminh2412/go-microservice-structure/service/app/integration/user/internal/pkg/error_helper"
)

func (s *Service) AccessCheckSystemUser(ctx context.Context, req *integration_user.SystemUserAccessCheckRequest) (*integration_user.SystemUserAccessCheckReply, error) {
	context_helper.AddTag(ctx, tag.ClientID, req.ClientID)
	context_helper.AddTag(ctx, tag.AccessUsername, req.AccessUsername)
	context_helper.AddTag(ctx, tag.TargetUsername, req.TargetUsername)

	p, err := extractSystemUserAccessCheckParams(req)
	if err != nil {
		return nil, err
	}
	sysUser := system_user.New(ctx, p.ClientID, p.TargetUsername)
	err = sysUser.AccessCheck(&system_user.AccessCheckRequest{
		AccessUsername: p.AccessUsername,
	})
	if err != nil {
		return nil, err
	}
	return &integration_user.SystemUserAccessCheckReply{}, nil
}

type systemUserAccessCheckParams struct {
	ClientID        string
	AccessUsername  string
	TargetUsername  string
	AllowSelfAccess bool
}

func extractSystemUserAccessCheckParams(req *integration_user.SystemUserAccessCheckRequest) (*systemUserAccessCheckParams, error) {
	if len(req.ClientID) == 0 {
		msg := "ClientID is required"
		return nil, error_helper.New(service_error.InvalidArgument, msg, nil)
	}
	if len(req.AccessUsername) == 0 {
		msg := "Access username is required"
		return nil, error_helper.New(service_error.InvalidArgument, msg, nil)
	}
	if len(req.TargetUsername) == 0 {
		msg := "Target username is required"
		return nil, error_helper.New(service_error.InvalidArgument, msg, nil)
	}

	p := &systemUserAccessCheckParams{
		ClientID:       req.ClientID,
		AccessUsername: req.AccessUsername,
		TargetUsername: req.TargetUsername,
	}
	return p, nil
}
