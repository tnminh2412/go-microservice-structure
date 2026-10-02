package handler

import (
	"context"
	"go-pkg/helper/context_helper"
	service_error "go-pkg/helper/error_helper"
	"proto-repo/generated/cmm/service/cmm/core/communication/core_communication"
	"proto-repo/generated/cmm/service/cmm/core/user/core_user"
	"proto-repo/generated/cmm/service/cmm/integration/user/integration_user"

	"github.com/tnminh2412/go-microservice-structure/service/app/integration/pkg/otp"
	"github.com/tnminh2412/go-microservice-structure/service/app/integration/user/internal/handler/defined/tag"
	"github.com/tnminh2412/go-microservice-structure/service/app/integration/user/internal/pkg/communication"
	"github.com/tnminh2412/go-microservice-structure/service/app/integration/user/internal/pkg/error_helper"
)

func (s *Service) ConfirmOtpReset(ctx context.Context, req *integration_user.OtpResetConfirmRequest) (*integration_user.OtpResetConfirmReply, error) {
	context_helper.AddTag(ctx, tag.ClientID, req.ClientID)

	if len(req.Token) == 0 {
		return nil, error_helper.New(service_error.InvalidArgument, "Token is required", nil)
	}
	id, err := s.encrypter.DecryptToken(req.Token)
	if err != nil {
		return nil, error_helper.New(integration_user.ErrorMessage_InvalidToken.String(), "", err)
	}

	rep, err := s.communicationSrv.Complete(ctx, &core_communication.CompleteRequest{
		ClientID:    req.ClientID,
		ID:          id,
		RequestTime: req.RequestTime,
		RequestBy:   req.RequestBy,
	})
	if err != nil {
		if _, ok := service_error.Is(err, core_communication.ErrorMessage_CommunicationExpired.String()); ok {
			return nil, error_helper.New(integration_user.ErrorMessage_InvalidTokenExpired.String(), "", err)
		}
		if _, ok := service_error.Is(err, core_communication.ErrorMessage_CommunicationCompleted.String()); ok {
			return nil, error_helper.New(integration_user.ErrorMessage_InvalidTokenCompleted.String(), "", err)
		}
		return nil, err
	}

	if rep.Resource != communication.Resource.OtpReset {
		return nil, error_helper.New(integration_user.ErrorMessage_InvalidToken.String(), "", nil)
	}

	_, err = s.userSrv.UpdateUser(ctx, &core_user.UpdateUserRequest{
		ClientID: req.ClientID,
		ID:       rep.ResourceID,
		Data: map[string]string{
			otp.OtpEnable: "false",
		},
		UpdatedTime: req.RequestTime,
		UpdatedBy:   req.RequestBy,
	})
	if err != nil {
		return nil, err
	}

	return &integration_user.OtpResetConfirmReply{}, nil
}
