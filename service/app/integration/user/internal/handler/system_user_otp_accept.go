package handler

import (
	"context"
	"go-pkg/helper/context_helper"
	service_error "go-pkg/helper/error_helper"
	"proto-repo/generated/cmm/service/cmm/integration/user/integration_user"
	"time"

	"github.com/tnminh2412/go-microservice-structure/service/app/integration/pkg/feature/person/user/system_user"
	"github.com/tnminh2412/go-microservice-structure/service/app/integration/user/internal/handler/defined/tag"
	"github.com/tnminh2412/go-microservice-structure/service/app/integration/user/internal/pkg/error_helper"

	"github.com/golang/protobuf/ptypes/timestamp"
)

func (s *Service) AcceptSystemUserOtp(ctx context.Context, req *integration_user.SystemUserAcceptOtpRequest) (*integration_user.SystemUserAcceptOtpReply, error) {
	context_helper.AddTag(ctx, tag.ClientID, req.ClientID)
	context_helper.AddTag(ctx, tag.Username, req.Username)

	p, err := s.extractSystemUserOtpAcceptParams(req)
	if err != nil {
		return nil, err
	}

	sys := system_user.New(ctx, p.ClientID, p.Username)
	err = sys.AcceptOtp(&system_user.OtpAcceptParams{
		Otp:         p.Otp,
		Time:        p.Time,
		UpdatedBy:   p.UpdatedBy,
		UpdatedTime: p.UpdatedTime.AsTime(),
		Token:       p.Token,
	})
	if err != nil {
		return nil, err
	}

	rep := &integration_user.SystemUserAcceptOtpReply{}
	return rep, nil
}

type systemUserOtpAcceptParams struct {
	ClientID    string
	Username    string
	Otp         string
	Time        time.Time
	UpdatedBy   string
	UpdatedTime *timestamp.Timestamp
	Token       string
}

func (s *Service) extractSystemUserOtpAcceptParams(req *integration_user.SystemUserAcceptOtpRequest) (*systemUserOtpAcceptParams, error) {
	if len(req.ClientID) == 0 {
		msg := "ClientID is required"
		return nil, error_helper.New(service_error.InvalidArgument, msg, nil)
	}
	if len(req.Username) == 0 {
		msg := "Username is required"
		return nil, error_helper.New(service_error.InvalidArgument, msg, nil)
	}
	if len(req.Otp) == 0 {
		msg := "Otp is required"
		return nil, error_helper.New(service_error.InvalidArgument, msg, nil)
	}
	if !req.Time.IsValid() {
		msg := "Time is invalid"
		return nil, error_helper.New(service_error.InvalidArgument, msg, nil)
	}

	p := &systemUserOtpAcceptParams{
		ClientID:    req.ClientID,
		Username:    req.Username,
		Otp:         req.Otp,
		Time:        req.Time.AsTime(),
		UpdatedBy:   req.UpdatedBy,
		UpdatedTime: req.UpdatedTime,
		Token:       req.Token,
	}
	return p, nil
}
