package adapter

import (
	"proto-repo/generated/cmm/service/cmm/core/user/core_user"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/tnminh2412/go-microservice-structure/service/app/core/user/internal/db"
)

func ToRoleReply(role *db.Role) *core_user.GetRoleReply {
	return &core_user.GetRoleReply{
		Role:  role.Name,
		Tasks: role.Tasks,
		Locale: func() []*core_user.Locale {
			var locale []*core_user.Locale

			for _, v := range role.Locale {
				locale = append(locale, &core_user.Locale{
					Language:    v.Language,
					Description: v.Description,
				})
			}

			return locale
		}(),
		ExcludeTasks: role.ExcludeTasks,
		CreatedBy:    role.CreatedBy,
		CreatedTime:  timestamppb.New(role.CreatedTime),
		UpdatedBy:    role.UpdatedBy,
		UpdatedTime:  timestamppb.New(role.UpdatedTime),
	}
}
