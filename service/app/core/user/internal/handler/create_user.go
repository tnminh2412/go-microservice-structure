package handler

import (
	"context"
	service_error "go-pkg/helper/error_helper"
	"proto-repo/generated/cmm/service/cmm/core/user/core_user"

	"github.com/tnminh2412/go-microservice-structure/service/app/core/user/internal/db"
	"github.com/tnminh2412/go-microservice-structure/service/app/core/user/internal/pkg/error_helper"

	"go.mongodb.org/mongo-driver/bson"
)

func (s *Service) CreateUser(ctx context.Context, req *core_user.CreateUserRequest) (*core_user.CreateUserReply, error) {
	mongo, err := s.SharedService.GetDatabase(ctx, req.ClientID)
	if err != nil {
		return nil, error_helper.New(service_error.Internal, "Server error", err)
	}

	userRequest, err := s.extractCreateUserRequest(req)
	if err != nil {
		return nil, err
	}
	err = db.CreateUser(ctx, mongo, userRequest)

	if err != nil {
		return nil, err
	}

	return &core_user.CreateUserReply{}, nil
}

func (s *Service) extractCreateUserRequest(req *core_user.CreateUserRequest) (*db.User, error) {

	contact, err := bson.Marshal(req.Contact)
	if err != nil {
		return nil, err
	}

	domains := getCreateDomains(req.Domains)

	newUser := &db.User{
		Username:      req.ID,
		FirstName:     req.FirstName,
		LastName:      req.LastName,
		DisplayName:   req.DisplayName,
		Roles:         req.Roles,
		Password:      req.Password,
		Status:        req.Status,
		ReferrerID:    req.ReferrerID,
		Contact:       contact,
		CurrencyCodes: req.CurrencyCodes,

		CreatedBy:   req.CreatedBy,
		CreatedTime: req.CreatedTime.AsTime(),
		UpdatedBy:   req.CreatedBy,
		UpdatedTime: req.CreatedTime.AsTime(),
		IsSystem:    req.IsSystem,

		Domains:                 domains,
		LandingRootDomains:      req.LandingRootDomains,
		IsRequireChangePassword: req.IsRequireChangePassword,
		IsEmailVerified:         req.IsEmailVerified,
		IsPhoneNumberVerified:   req.IsPhoneNumberVerified,

		Email:          req.Email,
		Phone:          req.Phone,
		PhoneCode:      req.PhoneCode,
		PlayerIDPrefix: req.PlayerIDPrefix,

		AccessLines: req.AccessLines,

		Level:        req.Level,
		DirectParent: req.ParentID,
		ParentL1:     req.ParentL1,
		ParentL2:     req.ParentL2,
		ParentL3:     req.ParentL3,
		ParentL4:     req.ParentL4,
		ParentL5:     req.ParentL5,
		ParentL6:     req.ParentL6,
		ParentL7:     req.ParentL7,
		ParentL8:     req.ParentL8,
		ParentL9:     req.ParentL9,
		ParentL10:    req.ParentL10,
		ParentL11:    req.ParentL11,
		ParentL12:    req.ParentL12,
		ParentL13:    req.ParentL13,
		ParentL14:    req.ParentL14,
		ParentL15:    req.ParentL15,
	}
	return newUser, nil
}

func getCreateDomains(domains []*core_user.Domain) []*db.Domain {
	dbDomains := make([]*db.Domain, len(domains))
	for i, d := range domains {
		dbDomains[i] = &db.Domain{
			Domain:    d.Domain,
			IsTest:    d.IsTest,
			Active:    d.Active,
			IsLobby:   d.IsLobby,
			Countries: d.Countries,
			Image:     d.Image,
			Name:      d.Name,
			Language:  d.Language,
			Providers: func() []*db.Provider {
				providers := make([]*db.Provider, len(d.Providers))
				for j, p := range d.Providers {
					providers[j] = &db.Provider{
						Provider: p.Provider,
						ID:       p.ID,
					}
				}
				return providers
			}(),
			CreatedBy: d.CreatedBy,
		}
	}

	return dbDomains
}
