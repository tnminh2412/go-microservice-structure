package handler

import (
	"context"
	"go-pkg/helper/context_helper"
	service_error "go-pkg/helper/error_helper"
	"proto-repo/generated/cmm/service/cmm/integration/user/integration_user"

	"github.com/tnminh2412/go-microservice-structure/service/app/integration/user/internal/adapter"
	"github.com/tnminh2412/go-microservice-structure/service/app/integration/user/internal/handler/defined/tag"
	"github.com/tnminh2412/go-microservice-structure/service/app/integration/user/internal/pkg/error_helper"
)

func (s *Service) GetValidationRule(ctx context.Context, req *integration_user.GetValidationRuleRequest) (*integration_user.GetValidationRuleReply, error) {
	if len(req.ClientID) == 0 {
		return nil, error_helper.New(service_error.InvalidArgument, "ClientID is required", nil)
	}

	context_helper.AddTag(ctx, tag.ClientID, req.ClientID)

	confLocal, err := s.LocalService.GetConfigValue(ctx, req.ClientID)
	if err != nil {
		return nil, err
	}

	rep := &integration_user.GetValidationRuleReply{
		UsernameRules: adapter.ToUserUsernameRules(confLocal.UsernameRules),
		PasswordRules: adapter.ToUserPasswordRules(confLocal.PasswordRules),
	}
	return rep, nil
}
