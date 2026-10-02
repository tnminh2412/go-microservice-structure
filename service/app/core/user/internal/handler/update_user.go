package handler

import (
	"context"
	service_error "go-pkg/helper/error_helper"
	"proto-repo/generated/cmm/service/cmm/core/user/core_user"

	"github.com/tnminh2412/go-microservice-structure/service/app/core/user/internal/db"
	"github.com/tnminh2412/go-microservice-structure/service/app/core/user/internal/pkg/error_helper"

	"go.mongodb.org/mongo-driver/bson"
)

func (s *Service) UpdateUser(ctx context.Context, req *core_user.UpdateUserRequest) (*core_user.UpdateUserReply, error) {
	mongo, err := s.SharedService.GetDatabase(ctx, req.ClientID)
	if err != nil {
		return nil, error_helper.New(service_error.Internal, "Server error", err)
	}

	contact, err := bson.Marshal(req.Contact)
	if err != nil {
		return nil, err
	}

	matched, err := db.UpdateUser(ctx, mongo, &db.UpdateUserParams{
		ByRoles:                   req.Roles != nil,
		ByTasks:                   req.Tasks != nil,
		ByPassword:                req.Password.IsValid() == true,
		ByStatus:                  req.Status.IsValid() == true,
		ByFirstName:               req.FirstName.IsValid() == true,
		ByLastName:                req.LastName.IsValid() == true,
		ByContact:                 req.Contact != nil,
		ByFailLoginAttempt:        req.FailLoginAttempt.IsValid() == true,
		ByCurrencyCodes:           req.CurrencyCodes != nil,
		ByReferrerID:              req.ReferrerID.IsValid() == true,
		ByDisplayName:             req.DisplayName.IsValid() == true,
		ByLastLoginTime:           req.LastLoginTime.IsValid() == true,
		ByLastLoginIP:             req.LastLoginIP.IsValid() == true,
		ByLastLoginCountry:        req.LastLoginCountry.IsValid() == true,
		ByIsRequireChangePassword: req.IsRequireChangePassword.IsValid() == true,
		ByData:                    len(req.Data) > 0,
		ByEmail:                   req.Email.IsValid(),
		ByDomains:                 req.Domains != nil,
		ByLandingRootDomains:      req.LandingRootDomains != nil,
		ByPhone:                   req.Phone.IsValid(),
		ByPhoneCode:               req.PhoneCode.IsValid(),
		ByIsPhoneNumberVerified:   req.IsPhoneNumberVerified.IsValid(),
		ByIsEmailVerified:         req.IsEmailVerified.IsValid(),
		ByImages:                  req.Images != nil,
		ByPlayerIDPrefix:          req.PlayerIDPrefix.IsValid(),

		User: &db.User{
			Username: req.ID,
			Roles: func() []string {
				if req.GetRoles() != nil {
					return req.GetRoles().List
				}
				return []string{}
			}(),
			Tasks: func() []string {
				if req.GetTasks() != nil {
					return req.GetTasks().GetList()
				}
				return []string{}
			}(),
			Password:         req.Password.GetValue(),
			Status:           req.Status.GetValue(),
			FirstName:        req.FirstName.GetValue(),
			LastName:         req.LastName.GetValue(),
			ReferrerID:       req.ReferrerID.GetValue(),
			DisplayName:      req.DisplayName.GetValue(),
			Contact:          contact,
			FailLoginAttempt: req.FailLoginAttempt.GetValue(),
			CurrencyCodes: func() []string {
				if req.GetCurrencyCodes() != nil {
					return req.GetCurrencyCodes().List
				}
				return []string{}
			}(),

			IsRequireChangePassword: req.IsRequireChangePassword.GetValue(),
			Data:                    req.Data,

			LastLoginCountry:      req.LastLoginCountry.GetValue(),
			LastLoginIP:           req.LastLoginIP.GetValue(),
			LastLoginTime:         req.LastLoginTime.AsTime(),
			UpdatedBy:             req.UpdatedBy,
			UpdatedTime:           req.UpdatedTime.AsTime(),
			Email:                 req.Email.GetValue(),
			Phone:                 req.Phone.GetValue(),
			PhoneCode:             req.PhoneCode.GetValue(),
			Domains:               getUpdateDomains(req.Domains),
			LandingRootDomains:    req.LandingRootDomains,
			IsEmailVerified:       req.IsEmailVerified.GetValue(),
			IsPhoneNumberVerified: req.IsPhoneNumberVerified.GetValue(),
			Images:                getUpdateImages(req.Images),
			PlayerIDPrefix:        req.PlayerIDPrefix.GetValue(),
		},
	})
	if err != nil {
		return nil, err
	}
	if matched == 0 {
		return nil, error_helper.New(service_error.NotFound, "Not found", err)
	}

	return &core_user.UpdateUserReply{}, nil
}

func getUpdateDomains(domains *core_user.Domains) []*db.Domain {
	if domains == nil {
		return nil
	}

	dbDomains := make([]*db.Domain, len(domains.List))
	for i, domain := range domains.List {
		dbDomain := &db.Domain{
			Domain:    domain.Domain,
			IsTest:    domain.IsTest,
			Active:    domain.Active,
			IsLobby:   domain.IsLobby,
			Countries: domain.Countries,
			Image:     domain.Image,
			Name:      domain.Name,
			Language:  domain.Language,
			Providers: func() []*db.Provider {
				providers := make([]*db.Provider, 0)
				for _, p := range domain.Providers {
					providers = append(providers, &db.Provider{
						Provider: p.Provider,
						ID:       p.ID,
					})
				}
				return providers
			}(),
			CreatedBy: domain.CreatedBy,
		}
		dbDomains[i] = dbDomain
	}

	return dbDomains
}

func getUpdateImages(images *core_user.Images) []*db.Image {
	if images == nil {
		return nil
	}

	dbImages := make([]*db.Image, len(images.List))
	for i, domain := range images.List {
		dbDomain := &db.Image{
			Image:    domain.Image,
			Name:     domain.Name,
			Language: domain.Language,
		}
		dbImages[i] = dbDomain
	}

	return dbImages
}
