package adapter

import (
	"encoding/json"
	"proto-repo/generated/cmm/service/cmm/core/user/core_user"

	"github.com/tnminh2412/go-microservice-structure/service/app/core/user/internal/db"

	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type RoleTasks struct {
	Tasks        []string
	ExcludeTasks []string
}

func ToUserReply(dbUser *db.User, roles []*db.Role) *core_user.GetUserReply {
	var mapRoles = make(map[string]*RoleTasks)

	for _, r := range roles {
		mapRoles[r.Name] = &RoleTasks{
			Tasks:        r.Tasks,
			ExcludeTasks: r.ExcludeTasks,
		}
	}

	var contact structpb.Struct
	_ = json.Unmarshal([]byte(dbUser.Contact.String()), &contact)

	rep := &core_user.GetUserReply{
		ID: dbUser.Username,
		Roles: func() []*core_user.TaskRole {
			var taskRoles []*core_user.TaskRole

			for _, r := range dbUser.Roles {
				taskRoles = append(taskRoles, &core_user.TaskRole{
					Role:         r,
					Tasks:        mapRoles[r].Tasks,
					Locales:      mapToLocales(r, roles),
					ExcludeTasks: mapRoles[r].ExcludeTasks,
				})
			}

			return taskRoles
		}(),
		AccessLines: dbUser.AccessLines,

		Tasks:         dbUser.Tasks,
		Password:      dbUser.Password,
		Status:        dbUser.Status,
		FirstName:     dbUser.FirstName,
		LastName:      dbUser.LastName,
		Level:         dbUser.Level,
		Contact:       &contact,
		CurrencyCodes: dbUser.CurrencyCodes,
		ReferrerID:    dbUser.ReferrerID,
		DisplayName:   dbUser.DisplayName,

		IsRequireChangePassword: dbUser.IsRequireChangePassword,
		IsEmailVerified:         dbUser.IsEmailVerified,
		IsPhoneNumberVerified:   dbUser.IsPhoneNumberVerified,

		LastLoginTime: func() *timestamppb.Timestamp {
			if dbUser.LastLoginTime.IsZero() {
				return nil
			}
			return timestamppb.New(dbUser.LastLoginTime)
		}(),
		LastLoginIP:      dbUser.LastLoginIP,
		LastLoginCountry: dbUser.LastLoginCountry,
		DirectParent:     dbUser.DirectParent,
		ParentL1:         dbUser.ParentL1,
		ParentL2:         dbUser.ParentL2,
		ParentL3:         dbUser.ParentL3,
		ParentL4:         dbUser.ParentL4,
		ParentL5:         dbUser.ParentL5,
		ParentL6:         dbUser.ParentL6,
		ParentL7:         dbUser.ParentL7,
		ParentL8:         dbUser.ParentL8,
		ParentL9:         dbUser.ParentL9,
		ParentL10:        dbUser.ParentL10,
		ParentL11:        dbUser.ParentL11,
		ParentL12:        dbUser.ParentL12,
		ParentL13:        dbUser.ParentL13,
		ParentL14:        dbUser.ParentL14,
		ParentL15:        dbUser.ParentL15,
		Path:             toPath(dbUser),

		IsExternal:         dbUser.IsExternal,
		IsSystem:           dbUser.IsSystem,
		IsTest:             dbUser.IsTest,
		FailLoginAttempt:   dbUser.FailLoginAttempt,
		CreatedBy:          dbUser.CreatedBy,
		CreatedTime:        timestamppb.New(dbUser.CreatedTime),
		UpdatedBy:          dbUser.UpdatedBy,
		UpdatedTime:        timestamppb.New(dbUser.UpdatedTime),
		Data:               dbUser.Data,
		Email:              dbUser.Email,
		Phone:              dbUser.Phone,
		PhoneCode:          dbUser.PhoneCode,
		Domains:            toDomainsReply(dbUser.Domains),
		LandingRootDomains: dbUser.LandingRootDomains,
		Images:             toImagesReply(dbUser.Images),
		PlayerIDPrefix:     dbUser.PlayerIDPrefix,
	}

	return rep
}

func mapToLocales(r string, roles []*db.Role) []*core_user.Locale {
	reply := make([]*core_user.Locale, 0)
	for _, role := range roles {
		if r == role.Name {
			for _, l := range role.Locale {
				reply = append(reply, &core_user.Locale{
					Language:    l.Language,
					Description: l.Description,
				})
			}
			break
		}
	}
	return reply
}

func toPath(dbUser *db.User) *core_user.Path {
	path := &core_user.Path{
		List: make([]*core_user.Path_User, 0),
	}

	defer func() {
		path.List = append(path.List, &core_user.Path_User{
			ID:    dbUser.Username,
			Level: dbUser.Level,
		})
	}()

	if len(dbUser.ParentL1) == 0 {
		return path
	}
	path.List = append(path.List, &core_user.Path_User{
		ID:    dbUser.ParentL1,
		Level: 1,
	})

	if len(dbUser.ParentL2) == 0 {
		return path
	}
	path.List = append(path.List, &core_user.Path_User{
		ID:    dbUser.ParentL2,
		Level: 2,
	})

	if len(dbUser.ParentL3) == 0 {
		return path
	}
	path.List = append(path.List, &core_user.Path_User{
		ID:    dbUser.ParentL3,
		Level: 3,
	})

	if len(dbUser.ParentL4) == 0 {
		return path
	}
	path.List = append(path.List, &core_user.Path_User{
		ID:    dbUser.ParentL4,
		Level: 4,
	})

	if len(dbUser.ParentL5) == 0 {
		return path
	}
	path.List = append(path.List, &core_user.Path_User{
		ID:    dbUser.ParentL5,
		Level: 5,
	})

	if len(dbUser.ParentL6) == 0 {
		return path
	}
	path.List = append(path.List, &core_user.Path_User{
		ID:    dbUser.ParentL6,
		Level: 6,
	})

	if len(dbUser.ParentL7) == 0 {
		return path
	}
	path.List = append(path.List, &core_user.Path_User{
		ID:    dbUser.ParentL7,
		Level: 7,
	})

	if len(dbUser.ParentL8) == 0 {
		return path
	}
	path.List = append(path.List, &core_user.Path_User{
		ID:    dbUser.ParentL8,
		Level: 8,
	})

	if len(dbUser.ParentL9) == 0 {
		return path
	}
	path.List = append(path.List, &core_user.Path_User{
		ID:    dbUser.ParentL9,
		Level: 9,
	})

	if len(dbUser.ParentL10) == 0 {
		return path
	}
	path.List = append(path.List, &core_user.Path_User{
		ID:    dbUser.ParentL10,
		Level: 10,
	})

	if len(dbUser.ParentL11) == 0 {
		return path
	}
	path.List = append(path.List, &core_user.Path_User{
		ID:    dbUser.ParentL11,
		Level: 11,
	})

	if len(dbUser.ParentL12) == 0 {
		return path
	}
	path.List = append(path.List, &core_user.Path_User{
		ID:    dbUser.ParentL12,
		Level: 12,
	})

	if len(dbUser.ParentL13) == 0 {
		return path
	}
	path.List = append(path.List, &core_user.Path_User{
		ID:    dbUser.ParentL13,
		Level: 13,
	})

	if len(dbUser.ParentL14) == 0 {
		return path
	}
	path.List = append(path.List, &core_user.Path_User{
		ID:    dbUser.ParentL14,
		Level: 14,
	})

	if len(dbUser.ParentL15) == 0 {
		return path
	}
	path.List = append(path.List, &core_user.Path_User{
		ID:    dbUser.ParentL15,
		Level: 15,
	})

	return path
}

func toDomainsReply(dbDomains []*db.Domain) []*core_user.Domain {
	domains := make([]*core_user.Domain, 0)

	for _, dbDomain := range dbDomains {
		domain := &core_user.Domain{
			Domain:    dbDomain.Domain,
			IsTest:    dbDomain.IsTest,
			Active:    dbDomain.Active,
			IsLobby:   dbDomain.IsLobby,
			Countries: dbDomain.Countries,
			Image:     dbDomain.Image,
			Name:      dbDomain.Name,
			Language:  dbDomain.Language,
			Providers: func() []*core_user.Provider {
				providers := make([]*core_user.Provider, 0)
				for _, p := range dbDomain.Providers {
					providers = append(providers, &core_user.Provider{
						Provider: p.Provider,
						ID:       p.ID,
					})
				}
				return providers
			}(),
			CreatedBy: dbDomain.CreatedBy,
		}

		domains = append(domains, domain)
	}

	return domains
}

func toImagesReply(dbImages []*db.Image) []*core_user.Image {
	images := make([]*core_user.Image, 0)

	for _, dbImage := range dbImages {
		image := &core_user.Image{
			Image:    dbImage.Image,
			Name:     dbImage.Name,
			Language: dbImage.Language,
		}
		images = append(images, image)
	}
	return images
}
