package handler

import (
	"context"
	"go-pkg/helper/context_helper"
	service_error "go-pkg/helper/error_helper"
	"proto-repo/generated/cmm/service/cmm/integration/user/integration_user"

	"github.com/tnminh2412/go-microservice-structure/service/app/integration/pkg/feature/person/user"
	"github.com/tnminh2412/go-microservice-structure/service/app/integration/user/internal/handler/defined/tag"
	"github.com/tnminh2412/go-microservice-structure/service/app/integration/user/internal/pkg/error_helper"
)

func (s *Service) GetLoginTokenInfo(ctx context.Context, in *integration_user.GetLoginTokenInfoRequest) (*integration_user.GetLoginTokenInfoReply, error) {
	context_helper.AddTag(ctx, tag.ClientID, in.ClientID)

	err := s.validateGetLoginTokenInfo(ctx, in)
	if err != nil {
		return nil, err
	}

	tokenData, err := user.GetToken(ctx, in.ClientID, in.Token)
	if err != nil {
		return nil, error_helper.New(service_error.Unauthenticated, "invalid token", err)
	}

	rep := &integration_user.GetLoginTokenInfoReply{
		Username:    tokenData.Username,
		ClientID:    tokenData.ClientID,
		Roles:       tokenData.Roles,
		Tasks:       tokenData.Tasks,
		AccessLines: tokenData.AccessLines,
		IsSystem:    tokenData.IsSystem,
		ParentL1:    tokenData.ParentL1,
		ParentL2:    tokenData.ParentL2,
		ParentL3:    tokenData.ParentL3,
		ParentL4:    tokenData.ParentL4,
		ParentL5:    tokenData.ParentL5,
		ParentL6:    tokenData.ParentL6,
		ParentL7:    tokenData.ParentL7,
		ParentL8:    tokenData.ParentL8,
		ParentL9:    tokenData.ParentL9,
		ParentL10:   tokenData.ParentL10,
		ParentL11:   tokenData.ParentL11,
		ParentL12:   tokenData.ParentL12,
		ParentL13:   tokenData.ParentL13,
		ParentL14:   tokenData.ParentL14,
		ParentL15:   tokenData.ParentL15,
		Level:       tokenData.Level,
		SubUsername: tokenData.SubUsername,
		IsMain:      tokenData.IsMain,
		RequireOtp:  tokenData.RequireOtp,
	}
	return rep, nil
}

func (s *Service) validateGetLoginTokenInfo(ctx context.Context, req *integration_user.GetLoginTokenInfoRequest) error {
	if len(req.Token) == 0 {
		msg := "Token is required"
		return error_helper.New(service_error.InvalidArgument, msg, nil)
	}
	if len(req.ClientID) == 0 {
		msg := "ClientID is required"
		return error_helper.New(service_error.InvalidArgument, msg, nil)
	}

	return nil
}
