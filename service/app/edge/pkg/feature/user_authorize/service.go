package user_authorize

import (
	"context"
	"encoding/json"
	"fmt"
	"go-pkg/conf"
	"go-pkg/grpc"
	"go-pkg/grpc/grpc_pkg/client_helper"
	"go-pkg/helper/context_helper"
	service_error "go-pkg/helper/error_helper"
	"go-pkg/helper/header_helper"
	"go-pkg/log"
	"net/url"
	"proto-repo/generated/cmm/service/cmm/integration/user/integration_user"

	"github.com/tnminh2412/go-microservice-structure/service/app/edge/pkg/feature/user_authorize/error_helper"
	"github.com/tnminh2412/go-microservice-structure/service/app/edge/pkg/feature/user_authorize/permission"
)

const (
	GroupAPI = "API"

	IntegrationUserHostFlag = "IntegrationUserService"
)

type TokenData struct {
	SubUsername    string
	Username       string
	Level          int32
	Token          string
	Roles          map[string]struct{}
	Tasks          map[string]struct{}
	AccessUsername string
	AccessLines    []string
	ParentL1       string
	ParentL2       string
	ParentL3       string
	ParentL4       string
	ParentL5       string
	ParentL6       string
	ParentL7       string
	ParentL8       string
	ParentL9       string
	ParentL10      string
	ParentL11      string
	ParentL12      string
	ParentL13      string
	ParentL14      string
	ParentL15      string
	IsSystem       bool
	IsMain         bool
	RequireOtp     bool
}

type IService interface {
	Get(ctx context.Context) (*TokenData, error)
	Authenticate(ctx context.Context, method string) (context.Context, error)
	IsAuthenticated(ctx context.Context) bool
}

type service struct {
	integrationUserService integration_user.ServiceClient
	anonymousMethods       map[string]struct{}
	skipOtpMethods         map[string]struct{}
	requireTaskMethods     map[string]string
	globalMethods          map[string]struct{}
}

func New(anonymousMethods []string, requireTaskMethods map[string]string, skipOtpMethods []string, globalMethods []string) IService {
	anonymous := map[string]struct{}{}
	for _, m := range anonymousMethods {
		anonymous[m] = struct{}{}
	}
	skipOtp := map[string]struct{}{}
	for _, m := range skipOtpMethods {
		skipOtp[m] = struct{}{}
	}
	_globalMethods := map[string]struct{}{}
	for _, m := range globalMethods {
		_globalMethods[m] = struct{}{}
	}
	integrationUserHost, _ := conf.GetString(GroupAPI, IntegrationUserHostFlag, "")
	integrationUserConn, err := grpc.NewClient(integrationUserHost)
	if err != nil {
		log.Error(func() string { return fmt.Sprintf("err: %s", err.Error()) })
		panic(err)
	}
	integrationUserService := integration_user.NewServiceClient(integrationUserConn)

	return &service{
		integrationUserService: integrationUserService,
		anonymousMethods:       anonymous,
		requireTaskMethods:     requireTaskMethods,
		skipOtpMethods:         skipOtp,
		globalMethods:          _globalMethods,
	}
}

func (s *service) Authenticate(ctx context.Context, method string) (context.Context, error) {
	_, isGlobalMethod := s.globalMethods[method]
	if isGlobalMethod {
		return ctx, nil
	}

	clientID, err := s.getClientID(ctx)
	if err != nil {
		return nil, err
	}
	if _, ok := s.anonymousMethods[method]; ok {
		return ctx, nil
	}
	token, err := s.getToken(ctx)
	if err != nil {
		return ctx, err
	}

	tokenData, err := s.getTokenData(ctx, clientID, token)
	if err != nil {
		return ctx, error_helper.New(service_error.Unauthenticated, "Token is not found", nil)
	}

	if s.isMatchUsername(ctx, tokenData.Username) == false {
		return ctx, error_helper.New(service_error.Unauthenticated, "Username not match", nil)
	}

	_, skipOtpMethod := s.skipOtpMethods[method]
	if !skipOtpMethod && tokenData.RequireOtp {
		return ctx, error_helper.New(service_error.PermissionDenied, "Otp is required", nil)
	}

	requireTask, ok := s.getRequireTask(method)
	if !ok {
		return ctx, error_helper.New(service_error.PermissionDenied, "invalid permission", nil)
	}
	if requireTask != permission.Task.None {
		if _, ok = tokenData.Tasks[permission.Task.All]; ok == false {
			if _, ok = tokenData.Tasks[requireTask]; ok == false {
				return ctx, error_helper.New(service_error.PermissionDenied, "invalid permission", nil)
			}
		}
	}

	tokenDataJson, err := json.Marshal(tokenData)
	if err != nil {
		return ctx, error_helper.New(service_error.Unauthenticated, "json.Marshal error", nil)
	}

	h := header_helper.New()
	h.Append("tokenData", string(tokenDataJson))

	ctx = context_helper.AppendClientContext(ctx, h)
	return ctx, nil
}

func (s *service) Get(ctx context.Context) (*TokenData, error) {
	headers, _ := context_helper.ClientRequest(ctx)
	tokenDataJson := headers.Get("tokenData")
	if len(tokenDataJson) > 0 {
		var token *TokenData
		err := json.Unmarshal([]byte(tokenDataJson[0]), &token)
		if err != nil {
			return nil, error_helper.New(service_error.Unauthenticated, "Token is not found", nil)
		}

		return token, nil
	}

	clientID, err := s.getClientID(ctx)
	if err != nil {
		return nil, error_helper.New(service_error.Unauthenticated, "Token is not found", nil)
	}
	t, err := s.getToken(ctx)
	if err != nil {
		return nil, error_helper.New(service_error.Unauthenticated, "Token is not found", nil)
	}

	token, err := s.getTokenData(ctx, clientID, t)
	if err != nil {
		return nil, error_helper.New(service_error.Unauthenticated, "Token is not found", nil)
	}

	return token, nil
}

func (s *service) IsAuthenticated(ctx context.Context) bool {
	_, err := s.Get(ctx)
	return err == nil
}

func (s *service) getTokenData(ctx context.Context, clientID string, token string) (*TokenData, error) {
	res, err := s.integrationUserService.GetLoginTokenInfo(ctx, &integration_user.GetLoginTokenInfoRequest{
		ClientID: clientID,
		Token:    token,
	})
	if err != nil {
		return nil, err
	}
	rep := &TokenData{
		SubUsername: res.SubUsername,
		Username:    res.Username,
		Level:       res.Level,
		Token:       token,
		Roles: func() map[string]struct{} {
			roles := make(map[string]struct{})
			for _, r := range res.Roles {
				roles[r] = struct{}{}
			}
			return roles
		}(),
		Tasks: func() map[string]struct{} {
			tasks := make(map[string]struct{})
			for _, t := range res.Tasks {
				tasks[t] = struct{}{}
			}
			return tasks
		}(),
		AccessUsername: func() string {
			if res.IsSystem == false {
				return res.Username
			}

			if len(res.AccessLines) > 0 {
				return res.AccessLines[0]
			}

			return ""
		}(),
		AccessLines: res.AccessLines,
		ParentL1:    res.ParentL1,
		ParentL2:    res.ParentL2,
		ParentL3:    res.ParentL3,
		ParentL4:    res.ParentL4,
		ParentL5:    res.ParentL5,
		ParentL6:    res.ParentL6,
		ParentL7:    res.ParentL7,
		ParentL8:    res.ParentL8,
		ParentL9:    res.ParentL9,
		ParentL10:   res.ParentL10,
		ParentL11:   res.ParentL11,
		ParentL12:   res.ParentL12,
		ParentL13:   res.ParentL13,
		ParentL14:   res.ParentL14,
		ParentL15:   res.ParentL15,
		IsSystem:    res.IsSystem,
		IsMain:      res.IsMain,
		RequireOtp:  res.RequireOtp,
	}

	return rep, nil
}

func (s *service) getClientID(ctx context.Context) (string, error) {
	values := context_helper.Request(ctx).InHeader("id")
	if values == nil || len(values) == 0 {
		return "", error_helper.New(service_error.InvalidArgument, "ClientID is not found", nil)
	}
	return values[0], nil
}

func (s *service) getToken(ctx context.Context) (string, error) {
	t := context_helper.Request(ctx).InHeader("t")
	if len(t) == 0 {
		return "", error_helper.New(service_error.Unauthenticated, "Token is not found", nil)
	}

	return t[0], nil
}

func (s *service) getRequireTask(method string) (string, bool) {
	per, ok := s.requireTaskMethods[method]
	return per, ok
}

func (s *service) isMatchUsername(ctx context.Context, username string) bool {
	requestQuery := client_helper.ClientRequestUrl(ctx)
	q, _ := url.ParseQuery(requestQuery)
	usernameQuery := q.Get("u")

	isMatch := usernameQuery == username
	return isMatch
}
