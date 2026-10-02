package handler

import (
	"context"
	"go-pkg/helper/context_helper"
	service_error "go-pkg/helper/error_helper"
	"go-pkg/helper/id_helper"
	"go-pkg/ocode"
	"net/url"
	"proto-repo/generated/cmm/service/cmm/core/communication/core_communication"
	"proto-repo/generated/cmm/service/cmm/core/user/core_user"
	"proto-repo/generated/cmm/service/cmm/integration/user/integration_user"
	"proto-repo/generated/pkg/nullable"
	"strings"
	"time"

	"github.com/tnminh2412/go-microservice-structure/service/app/integration/pkg/otp"
	"github.com/tnminh2412/go-microservice-structure/service/app/integration/user/internal/handler/defined/tag"
	"github.com/tnminh2412/go-microservice-structure/service/app/integration/user/internal/pkg/communication"
	"github.com/tnminh2412/go-microservice-structure/service/app/integration/user/internal/pkg/error_helper"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Service) RequestOtpReset(ctx context.Context, req *integration_user.OtpResetRequest) (*integration_user.OtpResetReply, error) {
	context_helper.AddTag(ctx, tag.ClientID, req.ClientID)

	email, err := s.validator.Email(req.Email)
	if err != nil {
		return nil, err
	}

	cfg, err := s.LocalService.GetConfigValue(ctx, req.ClientID)
	if err != nil {
		return nil, err
	}

	user, err := s.userSrv.GetUser(ctx, &core_user.GetUserRequest{
		ClientID: req.ClientID,
		Email:    &nullable.NullString{Value: email},
	})
	if err != nil {
		_, ok := service_error.Is(err, service_error.NotFound)
		if ok {
			msg := "Email is invalid"
			return nil, error_helper.New(integration_user.ErrorMessage_InvalidEmail.String(), msg, err)
		}
		return nil, err
	}

	otpEnable := otp.IsOtpEnable(user.Data)
	if !otpEnable {
		msg := "Otp is not enable"
		return nil, error_helper.New(service_error.InvalidArgument, msg, err)
	}

	id := ocode.String(id_helper.New())
	token, err := s.encrypter.GenerateToken(id)
	if err != nil {
		return nil, err
	}

	redirectUrl := strings.ReplaceAll(req.RedirectUrl, cfg.Otp.ResetTokenParam, url.QueryEscape(token))
	now := time.Now()
	expire := now.Add(time.Duration(cfg.Otp.ResetTokenTtlMinutes) * time.Minute)
	resend := now.Add(time.Duration(cfg.Otp.ResetTokenResendSeconds) * time.Second)
	_, err = s.communicationSrv.Send(ctx, &core_communication.SendRequest{
		ClientID:    req.ClientID,
		ID:          id,
		Resource:    communication.Resource.OtpReset,
		ResourceID:  user.ID,
		Destination: user.Email,
		SendType:    communication.SendType.Email,
		Template:    communication.Template.OtpReset,
		Data: map[string]string{
			communication.Data.Token:       token,
			communication.Data.Username:    user.ID,
			communication.Data.RedirectURI: redirectUrl,
		},
		Language:      req.Language,
		ValidFrom:     timestamppb.New(now),
		ExpiryAt:      timestamppb.New(expire),
		ValidResendAt: timestamppb.New(resend),
		RequestTime:   req.RequestTime,
		RequestBy:     req.RequestBy,
	})
	if err != nil {
		return nil, err
	}

	return &integration_user.OtpResetReply{}, nil
}
