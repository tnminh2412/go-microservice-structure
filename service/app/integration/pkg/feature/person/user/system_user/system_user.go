package system_user

import (
	"context"
	"go-pkg/helper/array_helper"
	service_error "go-pkg/helper/error_helper"
	"proto-repo/generated/cmm/service/cmm/core/user/core_user"
	"proto-repo/generated/pkg/nullable"
	"sort"
	"strconv"
	"time"

	"github.com/tnminh2412/go-microservice-structure/service/app/integration/pkg/feature/person/user/system_user/types/system_user_role"
	"github.com/tnminh2412/go-microservice-structure/service/app/integration/pkg/feature/person/user/system_user/types/system_user_status"
	"github.com/tnminh2412/go-microservice-structure/service/app/integration/pkg/feature/person/user/token_user"
	"github.com/tnminh2412/go-microservice-structure/service/app/integration/pkg/otp"
	"github.com/tnminh2412/go-microservice-structure/service/pkg/password"
	"github.com/tnminh2412/go-microservice-structure/service/pkg/validator"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type ISystemUser interface {
	Get() (*Info, error)
	Create(req *CreateRequest) error
	Update(req *UpdateRequest) error
	SetPassword(req *SetPasswordRequest) error
	ChangePassword(req *ChangePasswordRequest) error
	SetStatus(req *SetStatusRequest) error
	Login(req *LoginRequest) (*LoginResult, error)
	AccessCheck(req *AccessCheckRequest) error
	GenerateOtp(req *OtpGenerateParams) (*OtpGenerateResult, error)
	AcceptOtp(req *OtpAcceptParams) error
	ResetOtp(req *OtpResetParams) error
	SkipOtp(req *OtpSkipParams) error
}

var (
	userSrv   core_user.UserServiceClient
	tokenUser *token_user.Server
)

var (
	getUsernameRules       func(ctx context.Context, clientID string) ([]validator.UsernameRule, error)
	getPasswordRules       func(ctx context.Context, clientID string) ([]validator.PasswordRule, error)
	getUsernameBlacklist   func(ctx context.Context, clientID string) ([]string, error)
	getMaxFailLoginAttempt func(ctx context.Context, clientID string) (int32, error)
	getOtpIssuer           func(ctx context.Context, clientID string) (string, error)
	getChatConfig          func(ctx context.Context, clientID string) (*Chat, error)
	getTokenTTLConfig      func(ctx context.Context, clientID string) (time.Duration, error)
)

type SystemUser struct {
	info *Info

	ctx      context.Context
	clientID string
	username string
}

type ServiceInfo struct {
	UserService              core_user.UserServiceClient
	TokenUserService         *token_user.Server
	FnGetUsernameRules       func(ctx context.Context, clientID string) ([]validator.UsernameRule, error)
	FnGetPasswordRules       func(ctx context.Context, clientID string) ([]validator.PasswordRule, error)
	FnGetUsernameBlacklist   func(ctx context.Context, clientID string) ([]string, error)
	FnGetMaxFailLoginAttempt func(ctx context.Context, clientID string) (int32, error)
	FnGetOtpIssuer           func(ctx context.Context, clientID string) (string, error)
	FnGetChatConfig          func(ctx context.Context, clientID string) (*Chat, error)
	FnGetTokenTTLConfig      func(ctx context.Context, clientID string) (time.Duration, error)
}

func Init(info *ServiceInfo) {
	userSrv = info.UserService
	tokenUser = info.TokenUserService
	getUsernameRules = info.FnGetUsernameRules
	getPasswordRules = info.FnGetPasswordRules
	getUsernameBlacklist = info.FnGetUsernameBlacklist
	getMaxFailLoginAttempt = info.FnGetMaxFailLoginAttempt
	getOtpIssuer = info.FnGetOtpIssuer
	getChatConfig = info.FnGetChatConfig
	getTokenTTLConfig = info.FnGetTokenTTLConfig
}

func New(ctx context.Context, clientID string, username string) ISystemUser {
	return &SystemUser{
		ctx:      ctx,
		clientID: clientID,
		username: username,
	}
}

func (s *SystemUser) Get() (*Info, error) {
	request := &core_user.GetUserRequest{
		ClientID: s.clientID,
		ID:       &nullable.NullString{Value: s.username},
	}

	user, err := userSrv.GetUser(s.ctx, request)
	if err != nil {
		return nil, err
	}
	if !user.IsSystem {
		return nil, Error(service_error.NotFound, "", nil)
	}
	s.info = toSystemUser(user)
	if err != nil {
		return nil, err
	}
	chat, err := getChatInfo(s.ctx, s.clientID)
	if err != nil {
		return nil, err
	}
	s.info.Chat = chat

	return s.info, nil
}

func (s *SystemUser) Create(req *CreateRequest) error {
	p, err := s.extractCreateSystemUserRequest(req)
	if err != nil {
		return err
	}

	_, err = userSrv.CreateUser(s.ctx, &core_user.CreateUserRequest{
		ClientID:                p.ClientID,
		ID:                      p.Username,
		Roles:                   p.Roles,
		Password:                p.Password,
		Status:                  system_user_status.Active,
		FirstName:               p.FirstName,
		LastName:                p.LastName,
		DisplayName:             p.DisplayName,
		IsSystem:                true,
		IsRequireChangePassword: true,
		Level:                   0,
		Email:                   p.Email,
		AccessLines:             p.AccessLines,
		CreatedBy:               p.CreatedBy,
		CreatedTime:             timestamppb.New(p.CreatedTime),
	})
	if err != nil {
		if _, ok := service_error.Is(err, service_error.Duplicate); ok {
			return Error(InvalidUsernameDuplicated, "", err)
		}
		return err
	}
	return nil
}

func (s *SystemUser) Update(req *UpdateRequest) error {
	sysReq, err := s.extractUpdateSystemUserRequest(req)
	if err != nil {
		return err
	}
	updateRequest := &core_user.UpdateUserRequest{
		ClientID:    s.clientID,
		ID:          s.username,
		FirstName:   sysReq.FirstName,
		LastName:    sysReq.LastName,
		DisplayName: sysReq.DisplayName,
		Roles: func() *core_user.Roles {
			if sysReq.NewRoles != nil {
				return &core_user.Roles{List: sysReq.NewRoles}
			}
			return nil
		}(),
		UpdatedBy:   sysReq.UpdatedBy,
		UpdatedTime: timestamppb.New(sysReq.UpdatedTime),
		Data:        sysReq.Data,
		Email:       sysReq.Email,
	}
	_, err = userSrv.UpdateUser(s.ctx, updateRequest)
	if err != nil {
		return err
	}
	if sysReq.NewRoles != nil {
		touchRoles := !stringSliceEqual(sysReq.Roles, sysReq.NewRoles)
		if touchRoles {
			_ = tokenUser.RemoveByUsername(s.ctx, s.clientID, s.username)
		}
	}
	return nil
}

func (s *SystemUser) SetPassword(req *SetPasswordRequest) error {
	if len(req.Password) == 0 {
		return Error(service_error.InvalidArgument, "Password is required", nil)
	}

	user, err := userSrv.GetUser(s.ctx, &core_user.GetUserRequest{
		ClientID: s.clientID,
		ID:       &nullable.NullString{Value: s.username},
	})
	if err != nil {
		return err
	}
	if user.IsSystem == false {
		return Error(service_error.NotFound, "", nil)
	}

	pass, err := password.Hash(req.Password)
	if err != nil {
		return err
	}

	_, err = userSrv.UpdateUser(s.ctx, &core_user.UpdateUserRequest{
		ClientID:                s.clientID,
		ID:                      s.username,
		Password:                &nullable.NullString{Value: pass},
		IsRequireChangePassword: &nullable.NullBool{Value: true},
		UpdatedBy:               req.UpdatedBy,
		UpdatedTime:             timestamppb.New(req.UpdatedTime),
	})
	if err != nil {
		return err
	}

	_ = tokenUser.RemoveByUsername(s.ctx, s.clientID, s.username)
	return nil
}

func (s *SystemUser) ChangePassword(req *ChangePasswordRequest) error {
	usernameRules, err := getUsernameRules(s.ctx, s.clientID)
	if err != nil {
		return err
	}
	passwordRules, err := getPasswordRules(s.ctx, s.clientID)
	if err != nil {
		return err
	}

	usernameBlackList, err := getUsernameBlacklist(s.ctx, s.clientID)
	if err != nil {
		return err
	}
	validator := validator.New(passwordRules, usernameRules, usernameBlackList)
	if err = validator.Password(s.ctx, req.NewPassword, req.Language, "New password"); err != nil {
		return err
	}

	user, err := userSrv.GetUser(s.ctx, &core_user.GetUserRequest{
		ClientID: s.clientID,
		ID:       &nullable.NullString{Value: s.username},
	})
	if err != nil {
		return err
	}
	if user.IsSystem == false {
		return Error(service_error.NotFound, "", nil)
	}

	if err = password.Compare(req.OldPassword, user.Password); err != nil {
		return Error(InvalidCurrentPassword, "", err)
	}
	pass, err := password.Hash(req.NewPassword)
	if err != nil {
		return err
	}

	_, err = userSrv.UpdateUser(s.ctx, &core_user.UpdateUserRequest{
		ClientID:                s.clientID,
		ID:                      s.username,
		Password:                &nullable.NullString{Value: pass},
		IsRequireChangePassword: &nullable.NullBool{Value: false},
		UpdatedBy:               req.UpdatedBy,
		UpdatedTime:             timestamppb.New(req.UpdatedTime),
	})
	if err != nil {
		return err
	}

	return nil
}

var statuses = map[string]struct{}{
	system_user_status.Active: {},
	system_user_status.Locked: {},
}

func (s *SystemUser) SetStatus(req *SetStatusRequest) error {
	if _, ok := statuses[req.Status]; ok == false {
		return Error(service_error.InvalidArgument, "Invalid status", nil)
	}
	user, err := userSrv.GetUser(s.ctx, &core_user.GetUserRequest{
		ClientID: s.clientID,
		ID:       &nullable.NullString{Value: s.username},
	})
	if err != nil {
		return err
	}
	if user.IsSystem == false {
		return Error(service_error.NotFound, "", nil)
	}
	params := &core_user.UpdateUserRequest{
		ClientID:    s.clientID,
		ID:          s.username,
		Status:      &nullable.NullString{Value: req.Status},
		UpdatedBy:   req.UpdatedBy,
		UpdatedTime: timestamppb.New(req.UpdatedTime),
	}

	if req.Status == system_user_status.Active {
		params.FailLoginAttempt = &nullable.NullInt32{Value: 0}
	}

	_, err = userSrv.UpdateUser(s.ctx, params)
	if err != nil {
		return err
	}
	if req.Status != system_user_status.Active {
		_ = tokenUser.RemoveByUsername(s.ctx, s.clientID, s.username)
	}
	return nil
}

func (s *SystemUser) Login(req *LoginRequest) (*LoginResult, error) {
	user, err := userSrv.GetUser(s.ctx, &core_user.GetUserRequest{
		ClientID: s.clientID,
		ID:       &nullable.NullString{Value: s.username},
	})
	if err != nil {
		return nil, err
	}
	if user.IsSystem == false {
		return nil, Error(InvalidUsernameOrPassword, "", nil)
	}
	if user.Status == system_user_status.Locked {
		return nil, Error(InvalidStatusLocked, "", nil)
	}

	err = password.Compare(req.Password, user.Password)
	if err != nil {
		if _, ok := service_error.Is(err, password.InvalidUsernameOrPassword); ok {
			failAttempt := user.FailLoginAttempt + 1
			updateUser := &core_user.UpdateUserRequest{
				ClientID:         s.clientID,
				ID:               s.username,
				UpdatedBy:        req.UpdatedBy,
				UpdatedTime:      timestamppb.New(req.UpdatedTime),
				FailLoginAttempt: &nullable.NullInt32{Value: failAttempt},
			}

			maxFailLoginAttempt, _ := getMaxFailLoginAttempt(s.ctx, s.clientID)
			if failAttempt >= maxFailLoginAttempt {
				updateUser.Status = &nullable.NullString{Value: system_user_status.Locked}
			}
			_, _ = userSrv.UpdateUser(s.ctx, updateUser)
		}
		return nil, err
	}
	if len(user.Roles) == 0 {
		return nil, Error(service_error.PermissionDenied, "invalid roles", nil)
	}

	otpEnable := otp.IsOtpEnable(user.Data)
	if otpEnable {
		if len(req.Otp) == 0 {
			return nil, Error(InvalidOtpIsRequired, "", nil)
		}
		o, err := strconv.ParseInt(req.Otp, 10, 64)
		if err != nil {
			return nil, Error(InvalidOtp, "", err)
		}

		ok := validateOtp(user.Data, int(o), req.UpdatedTime)
		if !ok {
			return nil, Error(InvalidOtp, "", nil)
		}
	}

	roles := make([]string, 0)
	for _, r := range user.Roles {
		roles = append(roles, r.Role)
	}

	tokenTTL, err := getTokenTTLConfig(s.ctx, s.clientID)
	if err != nil {
		return nil, err
	}
	token, err := tokenUser.Generate(s.ctx, &token_user.SetData{
		SubUsername: s.username,
		Username:    s.username,
		ClientID:    s.clientID,
		Roles:       roles,
		Tasks:       getTasks(user.Roles),
		Level:       user.Level,
		AccessLines: user.AccessLines,
		IsSystem:    true,
		IsMain:      true,
		RequireOtp:  !otpEnable && !otp.IsOtpSkip(user.Data),
	}, tokenTTL)
	if err != nil {
		return nil, err
	}

	_, err = userSrv.UpdateUser(s.ctx, &core_user.UpdateUserRequest{
		ClientID:         s.clientID,
		ID:               s.username,
		LastLoginCountry: &nullable.NullString{Value: req.LastLoginCountry},
		LastLoginIP:      &nullable.NullString{Value: req.LastLoginIP},
		LastLoginTime:    timestamppb.New(req.LastLoginTime),
		FailLoginAttempt: &nullable.NullInt32{Value: 0},
		UpdatedBy:        req.UpdatedBy,
		UpdatedTime:      timestamppb.New(req.UpdatedTime),
	})
	if err != nil {
		return nil, err
	}

	return &LoginResult{
		Token:       token,
		IsOtpEnable: otpEnable,
		IsOtpSkip:   otp.IsOtpSkip(user.Data),
	}, nil
}

func (s *SystemUser) AccessCheck(req *AccessCheckRequest) error {
	return nil
}

func GetSystemRoles() []string {
	return []string{
		system_user_role.Admin,
		system_user_role.Root,
	}
}

func getTasks(taskRoles []*core_user.TaskRole) []string {
	tasks := make([]string, 0)
	for _, role := range taskRoles {
		for _, task := range role.Tasks {
			if array_helper.Contains(tasks, task, false) == false {
				if array_helper.Contains(role.ExcludeTasks, task, false) == false {
					tasks = append(tasks, task)
				}
			}
		}
	}
	return tasks
}

func getChatInfo(ctx context.Context, clientID string) (*Chat, error) {
	if getChatConfig == nil {
		return &Chat{}, nil
	}

	chatCfg, err := getChatConfig(ctx, clientID)
	if err != nil {
		return nil, err
	}

	return chatCfg, nil
}

func stringSliceEqual(first, second []string) bool {
	if len(first) != len(second) {
		return false
	}

	sort.Strings(first)
	sort.Strings(second)

	for i, v := range second {
		if v != first[i] {
			return false
		}
	}

	return true
}
