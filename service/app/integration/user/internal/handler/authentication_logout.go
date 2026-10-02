package handler

import (
	"context"
	service_error "go-pkg/helper/error_helper"
	"proto-repo/generated/cmm/service/cmm/integration/user/integration_user"

	"github.com/tnminh2412/go-microservice-structure/service/app/integration/pkg/feature/person/user"
	"github.com/tnminh2412/go-microservice-structure/service/app/integration/user/internal/pkg/error_helper"
)

func (s *Service) Logout(ctx context.Context, req *integration_user.LogoutRequest) (*integration_user.LogoutReply, error) {
	p, err := s.validateLogoutRequest(req)
	if err != nil {
		return nil, err
	}
	err = user.RemoveToken(ctx, p.ClientID, p.Token)
	if err != nil {
		return nil, error_helper.New(service_error.Internal, "", err)
	}
	return &integration_user.LogoutReply{}, nil
}

type logoutParams struct {
	ClientID string
	Token    string
}

func (s *Service) validateLogoutRequest(req *integration_user.LogoutRequest) (*logoutParams, error) {
	if len(req.ClientID) == 0 {
		msg := "ClientId is required"
		return nil, error_helper.New(service_error.InvalidArgument, msg, nil)
	}
	if len(req.Token) == 0 {
		msg := "Token is required"
		return nil, error_helper.New(service_error.InvalidArgument, msg, nil)
	}
	p := &logoutParams{
		Token:    req.Token,
		ClientID: req.ClientID,
	}
	return p, nil
}
