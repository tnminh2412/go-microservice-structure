package handler

import (
	"context"
	"go-pkg/helper/context_helper"
	service_error "go-pkg/helper/error_helper"
	"proto-repo/generated/cmm/service/cmm/core/user/core_user"
	"proto-repo/generated/cmm/service/cmm/integration/user/integration_user"
	"proto-repo/generated/pkg/nullable"

	"github.com/tnminh2412/go-microservice-structure/service/app/integration/pkg/feature/person/user/system_user"
	"github.com/tnminh2412/go-microservice-structure/service/app/integration/user/internal/adapter"
	"github.com/tnminh2412/go-microservice-structure/service/app/integration/user/internal/handler/defined/tag"
	"github.com/tnminh2412/go-microservice-structure/service/app/integration/user/internal/pkg/error_helper"
)

func (s *Service) ListSystemUser(ctx context.Context, req *integration_user.SystemUserListRequest) (*integration_user.SystemUserListReply, error) {
	context_helper.AddTag(ctx, tag.ClientID, req.ClientID)

	p, err := s.extractListSystemUserParams(req)
	if err != nil {
		return nil, err
	}

	params := &core_user.ListUserRequest{
		ClientID:      p.ClientID,
		Statuses:      p.Statuses,
		IsSystem:      &nullable.NullBool{Value: true},
		ExcludeRoles:  system_user.GetSystemRoles(),
		Skip:          &nullable.NullInt64{Value: p.Skip},
		Limit:         &nullable.NullInt64{Value: p.Limit},
		Roles:         p.Roles,
		ExcludeHidden: &nullable.NullBool{Value: true},
	}
	if len(p.Search) > 0 {
		params.IDs = []string{p.Search}
	}
	res, err := s.userSrv.ListUser(ctx, params)
	if err != nil {
		return nil, err
	}

	rep := &integration_user.SystemUserListReply{
		Users: adapter.ToSystemUsers(res.Users),
	}
	return rep, nil
}

type systemUserListParams struct {
	ClientID     string
	Search       string
	Statuses     []string
	Skip         int64
	Limit        int64
	ExcludeRoles []string
	Roles        []string
}

func (s *Service) extractListSystemUserParams(req *integration_user.SystemUserListRequest) (*systemUserListParams, error) {
	if len(req.ClientID) == 0 {
		msg := "ClientID is required"
		return nil, error_helper.New(service_error.InvalidArgument, msg, nil)
	}
	err := s.validator.Page(req.Skip, req.Limit)
	if err != nil {
		return nil, err
	}
	p := &systemUserListParams{
		ClientID: req.ClientID,
		Search:   req.Search,
		Statuses: req.Statuses,
		Skip:     req.Skip,
		Limit:    req.Limit,
		Roles:    req.Roles,
	}
	return p, nil
}
