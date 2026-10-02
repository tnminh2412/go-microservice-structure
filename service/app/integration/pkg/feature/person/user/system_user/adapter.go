package system_user

import (
	"proto-repo/generated/cmm/service/cmm/core/user/core_user"
)

func toRoles(roles []*core_user.TaskRole) []*TaskRole {
	reply := make([]*TaskRole, len(roles))
	for i, r := range roles {
		reply[i] = &TaskRole{
			Role:  r.Role,
			Tasks: r.Tasks,
			Locales: func() []*Locale {
				reply := make([]*Locale, 0)
				for _, v := range r.Locales {
					reply = append(reply, &Locale{
						Description: v.Description,
						Language:    v.Language,
					})
				}
				return reply
			}(),
		}
	}
	return reply
}

func toSystemUser(u *core_user.GetUserReply) *Info {
	return &Info{
		Username:                u.ID,
		FirstName:               u.FirstName,
		LastName:                u.LastName,
		DisplayName:             u.DisplayName,
		Roles:                   toRoles(u.Roles),
		Password:                u.Password,
		Level:                   u.Level,
		IsSystem:                u.IsSystem,
		LastLoginCountry:        u.LastLoginCountry,
		LastLoginIP:             u.LastLoginIP,
		LastLoginTime:           u.LastLoginTime.AsTime(),
		CreatedBy:               u.CreatedBy,
		CreatedTime:             u.CreatedTime.AsTime(),
		UpdatedBy:               u.UpdatedBy,
		UpdatedTime:             u.UpdatedTime.AsTime(),
		Status:                  u.Status,
		IsRequireChangePassword: u.IsRequireChangePassword,
		Data:                    u.Data,
		Email:                   u.Email,
	}
}
