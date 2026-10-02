package handler

import (
	"context"
	"go-pkg/helper/context_helper"
	service_error "go-pkg/helper/error_helper"
	"proto-repo/generated/cmm/service/cmm/integration/user/integration_user"

	"github.com/tnminh2412/go-microservice-structure/service/app/integration/pkg/feature/person/user/system_user"
	"github.com/tnminh2412/go-microservice-structure/service/app/integration/user/internal/handler/defined/tag"
	"github.com/tnminh2412/go-microservice-structure/service/app/integration/user/internal/pkg/error_helper"

	"github.com/golang/protobuf/ptypes/timestamp"
)

func (s *Service) ResetSystemUserOtp(ctx context.Context, req *integration_user.SystemUserResetOtpRequest) (*integration_user.SystemUserResetOtpReply, error) {
	context_helper.AddTag(ctx, tag.ClientID, req.ClientID)
	context_helper.AddTag(ctx, tag.Username, req.Username)

	p, err := s.extractSystemUserOtpResetParams(req)
	if err != nil {
		return nil, err
	}

	sys := system_user.New(ctx, p.ClientID, p.Username)
	err = sys.ResetOtp(&system_user.OtpResetParams{
		UpdatedBy:   p.UpdatedBy,
		UpdatedTime: p.UpdatedTime.AsTime(),
	})
	if err != nil {
		return nil, err
	}

	rep := &integration_user.SystemUserResetOtpReply{}
	return rep, nil
}

type systemUserOtpResetParams struct {
	ClientID    string
	Username    string
	UpdatedBy   string
	UpdatedTime *timestamp.Timestamp
}

func (s *Service) extractSystemUserOtpResetParams(req *integration_user.SystemUserResetOtpRequest) (*systemUserOtpResetParams, error) {
	if len(req.ClientID) == 0 {
		msg := "ClientID is required"
		return nil, error_helper.New(service_error.InvalidArgument, msg, nil)
	}
	if len(req.Username) == 0 {
		msg := "Username is required"
		return nil, error_helper.New(service_error.InvalidArgument, msg, nil)
	}

	p := &systemUserOtpResetParams{
		ClientID:    req.ClientID,
		Username:    req.Username,
		UpdatedBy:   req.UpdatedBy,
		UpdatedTime: req.UpdatedTime,
	}
	return p, nil
}
