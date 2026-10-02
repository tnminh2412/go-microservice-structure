package access_check

import (
	"context"
	"proto-repo/generated/cmm/service/cmm/integration/user/integration_user"

	"github.com/tnminh2412/go-microservice-structure/service/app/edge/backoffice/internal/pkg/local_conf"
	"github.com/tnminh2412/go-microservice-structure/service/app/edge/pkg/feature/user_authorize"
)

// Server verifies that the logged-in user is allowed to act on a target account.
type Server struct {
	integrationUserClient integration_user.ServiceClient
	localConf             local_conf.IService
	userAuthorize         user_authorize.IService
}

func New(integrationUserClient integration_user.ServiceClient,
	localConf local_conf.IService,
	userAuthorize user_authorize.IService) *Server {
	return &Server{
		integrationUserClient: integrationUserClient,
		localConf:             localConf,
		userAuthorize:         userAuthorize,
	}
}

type AccessSystemUser struct {
	ClientID       string
	TargetUsername string
}

func (s *Server) AccessSystemUser(ctx context.Context, accessSystemUser *AccessSystemUser) error {
	tokenData, err := s.userAuthorize.Get(ctx)
	if err != nil {
		return err
	}
	_, err = s.integrationUserClient.AccessCheckSystemUser(ctx, &integration_user.SystemUserAccessCheckRequest{
		ClientID:       accessSystemUser.ClientID,
		AccessUsername: tokenData.Username,
		TargetUsername: accessSystemUser.TargetUsername,
	})
	if err != nil {
		return err
	}
	return nil
}
