package adapter

import (
	"proto-repo/generated/cmm/service/cmm/edge/backoffice/edge_backoffice"
	"proto-repo/generated/cmm/service/cmm/integration/user/integration_user"
	"sort"
)

func ToSystemUser(u *integration_user.SystemUser) *edge_backoffice.SystemUser {
	return &edge_backoffice.SystemUser{
		Username: u.Username,
		Roles: func() []string {
			list := make([]string, 0)
			for _, r := range u.Roles {
				list = append(list, r.Role)
			}
			sort.Strings(list)
			return list
		}(),
		UserRoles: func() []*edge_backoffice.Role {
			list := make([]*edge_backoffice.Role, 0)
			for _, r := range u.Roles {
				list = append(list, &edge_backoffice.Role{
					Name: r.Role,
					Locales: func() []*edge_backoffice.Locale {
						lstLocale := make([]*edge_backoffice.Locale, 0)
						for _, locale := range r.Locales {
							lstLocale = append(lstLocale, &edge_backoffice.Locale{
								Language:    locale.Language,
								Description: locale.Description,
							})
						}

						return lstLocale
					}(),
				})
			}

			sort.Slice(list, func(i int, j int) bool {
				return list[i].Name > list[j].Name
			})

			return list
		}(),
		Status:           u.Status,
		FirstName:        u.FirstName,
		LastName:         u.LastName,
		CreatedBy:        u.CreatedBy,
		CreatedTime:      u.CreatedTime,
		UpdatedBy:        u.UpdatedBy,
		UpdatedTime:      u.UpdatedTime,
		LastLoginCountry: u.LastLoginCountry,
		LastLoginIP:      u.LastLoginIP,
		LastLoginTime:    u.LastLoginTime,
		Email:            u.Email,
	}
}

func ToSystemUsers(users []*integration_user.SystemUser) []*edge_backoffice.SystemUser {
	arr := make([]*edge_backoffice.SystemUser, len(users))
	for i, u := range users {
		arr[i] = ToSystemUser(u)
	}
	return arr
}

func toUserBalance(data *integration_user.UserBalance) *edge_backoffice.UserBalance {
	return &edge_backoffice.UserBalance{
		CurrencyCode: data.CurrencyCode,
		Balance:      data.Balance,
	}
}
