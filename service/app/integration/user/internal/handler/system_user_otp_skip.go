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

func (s *Service) SkipSystemUserOtp(ctx context.Context, req *integration_user.SystemUserSkipOtpRequest) (*integration_user.SystemUserSkipOtpReply, error) {
	context_helper.AddTag(ctx, tag.ClientID, req.ClientID)
	context_helper.AddTag(ctx, tag.Username, req.Username)

	p, err := s.extractSystemUserOtpSkipParams(req)
	if err != nil {
		return nil, err
	}

	sys := system_user.New(ctx, p.ClientID, p.Username)
	err = sys.SkipOtp(&system_user.OtpSkipParams{
		UpdatedBy:   p.UpdatedBy,
		UpdatedTime: p.UpdatedTime.AsTime(),
		Value:       p.Value,
	})
	if err != nil {
		return nil, err
	}

	rep := &integration_user.SystemUserSkipOtpReply{}
	return rep, nil
}

type systemUserOtpSkipParams struct {
	ClientID    string
	Username    string
	UpdatedBy   string
	UpdatedTime *timestamp.Timestamp
	Value       bool
}

func (s *Service) extractSystemUserOtpSkipParams(req *integration_user.SystemUserSkipOtpRequest) (*systemUserOtpSkipParams, error) {
	if len(req.ClientID) == 0 {
		msg := "ClientID is required"
		return nil, error_helper.New(service_error.InvalidArgument, msg, nil)
	}
	if len(req.Username) == 0 {
		msg := "Username is required"
		return nil, error_helper.New(service_error.InvalidArgument, msg, nil)
	}

	p := &systemUserOtpSkipParams{
		ClientID:    req.ClientID,
		Username:    req.Username,
		UpdatedBy:   req.UpdatedBy,
		UpdatedTime: req.UpdatedTime,
		Value:       req.Value,
	}
	return p, nil
}
