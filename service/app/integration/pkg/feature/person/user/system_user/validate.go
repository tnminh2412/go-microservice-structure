package system_user

import (
	service_error "go-pkg/helper/error_helper"
	"proto-repo/generated/cmm/service/cmm/core/user/core_user"
	"proto-repo/generated/pkg/nullable"

	"github.com/tnminh2412/go-microservice-structure/service/pkg/password"
	"github.com/tnminh2412/go-microservice-structure/service/pkg/validator"
)

func (s *SystemUser) extractCreateSystemUserRequest(req *CreateRequest) (*internalCreateRequest, error) {
	usernameRules, err := getUsernameRules(s.ctx, s.clientID)
	if err != nil {
		return nil, err
	}
	passwordRules, err := getPasswordRules(s.ctx, s.clientID)
	if err != nil {
		return nil, err
	}

	usernameBlackList, err := getUsernameBlacklist(s.ctx, s.clientID)
	if err != nil {
		return nil, err
	}
	validator := validator.New(passwordRules, usernameRules, usernameBlackList)
	if err = validator.Username(s.ctx, s.username, req.Language); err != nil {
		return nil, err
	}

	if len(req.Password) == 0 {
		return nil, Error(InvalidPasswordIsRequired, "Password is required", nil)
	}

	email := ""
	if len(req.Email) != 0 {
		m, err := validator.Email(req.Email)
		if err != nil {
			return nil, err
		}
		email = m
	}

	err = s.validateRoles(req.Roles)
	if err != nil {
		return nil, err
	}
	pass, err := password.Hash(req.Password)
	if err != nil {
		return nil, err
	}
	reply := &internalCreateRequest{
		ClientID:    s.clientID,
		Username:    s.username,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Password:    pass,
		Roles:       req.Roles,
		DisplayName: req.DisplayName,
		CreatedBy:   req.CreatedBy,
		CreatedTime: req.CreatedTime,
		Email:       email,
		AccessLines: req.AccessLines,
	}
	return reply, nil
}

func (s *SystemUser) extractUpdateSystemUserRequest(req *UpdateRequest) (*internalUpdateRequest, error) {
	err := s.validateRoles(req.Roles)
	if err != nil {
		return nil, err
	}

	user, err := userSrv.GetUser(s.ctx, &core_user.GetUserRequest{
		ClientID: s.clientID,
		ID:       &nullable.NullString{Value: s.username},
	})
	if err != nil {
		return nil, err
	}
	if !user.IsSystem {
		return nil, Error(service_error.NotFound, "", nil)
	}

	updateRequest := &internalUpdateRequest{
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		DisplayName: req.DisplayName,
		Roles: func() []string {
			roles := make([]string, 0)
			for _, r := range user.Roles {
				roles = append(roles, r.Role)
			}
			return roles
		}(),
		NewRoles:    req.Roles,
		UpdatedBy:   req.UpdatedBy,
		UpdatedTime: req.UpdatedTime,
		Data:        req.Data,
	}

	if req.Email.IsValid() {
		usernameRules, err := getUsernameRules(s.ctx, s.clientID)
		if err != nil {
			return nil, err
		}
		passwordRules, err := getPasswordRules(s.ctx, s.clientID)
		if err != nil {
			return nil, err
		}

		usernameBlackList, err := getUsernameBlacklist(s.ctx, s.clientID)
		if err != nil {
			return nil, err
		}
		validator := validator.New(passwordRules, usernameRules, usernameBlackList)
		m, err := validator.Email(req.Email.GetValue())
		if err != nil {
			return nil, err
		}
		updateRequest.Email = &nullable.NullString{Value: m}
	}

	return updateRequest, nil
}

func (s *SystemUser) validateRoles(roles []string) error {
	if len(roles) == 0 {
		return nil
	}
	res, err := userSrv.ListRole(s.ctx, &core_user.ListRoleRequest{
		ClientID:      s.clientID,
		ExcludeSystem: true,
	})

	if err != nil {
		return err
	}
	allowedRoles := make(map[string]struct{})
	for _, r := range res.Roles {
		allowedRoles[r.Role] = struct{}{}
	}
	for _, r := range roles {
		_, ok := allowedRoles[r]
		if !ok {
			msg := "Roles is invalid"
			return Error(service_error.InvalidArgument, msg, nil)
		}
	}
	return nil
}
