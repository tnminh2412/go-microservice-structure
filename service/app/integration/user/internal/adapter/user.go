package adapter

import (
	"proto-repo/generated/cmm/service/cmm/core/user/core_user"
	"proto-repo/generated/cmm/service/cmm/integration/user/integration_user"

	"github.com/tnminh2412/go-microservice-structure/service/app/integration/pkg/feature/person/user/system_user"
	"github.com/tnminh2412/go-microservice-structure/service/app/integration/pkg/otp"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToRoles(roles []*core_user.TaskRole) []*integration_user.TaskRole {
	arr := make([]*integration_user.TaskRole, len(roles))
	for i, r := range roles {
		arr[i] = ToRole(r)
	}
	return arr
}

func ToRole(r *core_user.TaskRole) *integration_user.TaskRole {

	return &integration_user.TaskRole{
		Role:  r.Role,
		Tasks: r.Tasks,
		Locales: func() []*integration_user.Locale {
			reply := make([]*integration_user.Locale, 0)
			for _, v := range r.Locales {
				reply = append(reply, &integration_user.Locale{
					Description: v.Description,
					Language:    v.Language,
				})
			}
			return reply
		}(),
	}
}

func ToGetSystemUser(u *system_user.Info) *integration_user.SystemUser {
	return &integration_user.SystemUser{
		Username: u.Username,
		Roles: func() []*integration_user.TaskRole {
			var roles []*integration_user.TaskRole
			for _, r := range u.Roles {
				roles = append(roles, &integration_user.TaskRole{
					Role:  r.Role,
					Tasks: r.Tasks,
					Locales: func() []*integration_user.Locale {
						reply := make([]*integration_user.Locale, 0)
						for _, v := range r.Locales {
							reply = append(reply, &integration_user.Locale{
								Description: v.Description,
								Language:    v.Language,
							})
						}

						return reply
					}(),
				})
			}
			return roles
		}(),
		FirstName:               u.FirstName,
		LastName:                u.LastName,
		DisplayName:             u.DisplayName,
		Status:                  u.Status,
		IsRequireChangePassword: u.IsRequireChangePassword,
		CreatedBy:               u.CreatedBy,
		CreatedTime:             timestamppb.New(u.CreatedTime),
		UpdatedBy:               u.UpdatedBy,
		UpdatedTime:             timestamppb.New(u.UpdatedTime),
		LastLoginCountry:        u.LastLoginCountry,
		LastLoginIP:             u.LastLoginIP,
		LastLoginTime:           timestamppb.New(u.LastLoginTime),
		Email:                   u.Email,
		IsOtpSkip:               otp.IsOtpSkip(u.Data),
		Chat:                    toSystemUserChatReply(u.Chat),
	}
}

func ToSystemUser(u *core_user.GetUserReply) *integration_user.SystemUser {
	return &integration_user.SystemUser{
		Username:         u.ID,
		Roles:            ToRoles(u.Roles),
		Status:           u.Status,
		FirstName:        u.FirstName,
		LastName:         u.LastName,
		DisplayName:      u.DisplayName,
		CreatedBy:        u.CreatedBy,
		CreatedTime:      u.CreatedTime,
		UpdatedBy:        u.UpdatedBy,
		UpdatedTime:      u.UpdatedTime,
		LastLoginCountry: u.LastLoginCountry,
		LastLoginIP:      u.LastLoginIP,
		LastLoginTime:    u.LastLoginTime,
		Email:            u.Email,
		IsOtpSkip:        otp.IsOtpSkip(u.Data),
	}
}

func ToSystemUsers(users []*core_user.GetUserReply) []*integration_user.SystemUser {
	arr := make([]*integration_user.SystemUser, len(users))
	for i, u := range users {
		arr[i] = ToSystemUser(u)
	}
	return arr
}

func toSystemUserChatReply(userChat *system_user.Chat) *integration_user.Chat {
	if userChat == nil {
		return nil
	}
	chat := &integration_user.Chat{
		FreshWork: &integration_user.FreshWork{
			PlayerApp: func() *integration_user.App {
				if userChat.FreshWork == nil {
					return nil
				}
				return toSystemUserApp(userChat.FreshWork.PlayerApp)
			}(),
			AgentID: func() string {
				if userChat.FreshWork == nil {
					return ""
				}
				return userChat.FreshWork.AgentID
			}(),
			Enable: func() bool {
				if userChat.FreshWork == nil {
					return false
				}
				return userChat.FreshWork.Enable
			}(),
		},
	}
	return chat
}

func toSystemUserApp(app *system_user.App) *integration_user.App {
	if app == nil {
		return nil
	}
	return &integration_user.App{
		Ticket:      toSystemFreshWorkApp(app.Ticket),
		ChatSupport: toSystemFreshWorkApp(app.ChatSupport),
	}
}

func toSystemFreshWorkApp(freshWorkApp *system_user.FreshWorkApp) *integration_user.FreshWorkApp {
	if freshWorkApp == nil {
		return nil
	}
	return &integration_user.FreshWorkApp{
		ClientID:     freshWorkApp.ClientID,
		LoginPageUri: freshWorkApp.LoginPageUri,
		Name:         freshWorkApp.Name,
		Description:  freshWorkApp.Description,
	}
}
