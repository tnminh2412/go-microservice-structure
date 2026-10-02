package handler

import (
	"context"
	service_error "go-pkg/helper/error_helper"
	"go-pkg/mongo/mongo_error"
	"proto-repo/generated/cmm/service/cmm/core/user/core_user"

	"github.com/tnminh2412/go-microservice-structure/service/app/core/user/internal/adapter"
	"github.com/tnminh2412/go-microservice-structure/service/app/core/user/internal/db"
	"github.com/tnminh2412/go-microservice-structure/service/app/core/user/internal/pkg/error_helper"
)

func (s *Service) GetUser(ctx context.Context, req *core_user.GetUserRequest) (*core_user.GetUserReply, error) {
	mongo, err := s.SharedService.GetDatabase(ctx, req.ClientID)
	if err != nil {
		return nil, error_helper.New(service_error.Internal, "Server error", err)
	}

	p := db.GetUserParams{
		ByID:         req.ID.IsValid() == true,
		ByReferrerID: req.ReferrerID.IsValid() == true,
		ByEmail:      req.Email.IsValid() == true,
		ByDomain:     req.Domain.IsValid() == true,
		ByPhone:      req.Phone.IsValid() == true,
		ByPhoneCode:  req.PhoneCode.IsValid() == true,

		ID:         req.ID.GetValue(),
		ReferrerID: req.ReferrerID.GetValue(),
		Email:      req.Email.GetValue(),
		Domain:     req.Domain.GetValue(),
		Phone:      req.Phone.GetValue(),
		PhoneCode:  req.PhoneCode.GetValue(),
	}

	fields := getDBExcludeMask(req.FieldMask)
	p.Mask = &db.Mask{
		Fields:    fields,
		IsExclude: true,
	}

	user, err := db.GetUser(ctx, mongo, p)
	if err != nil {
		if mongo_error.IsNotFound(err) {
			return nil, error_helper.New(service_error.NotFound, "Not found", err)
		}

		return nil, error_helper.New(service_error.Internal, "Server error", err)
	}

	roles, err := db.ListRole(ctx, mongo, db.ListRoleParams{})
	if err != nil {
		return nil, error_helper.New(service_error.Internal, "Server error", err)
	}

	return adapter.ToUserReply(user, roles), nil
}

var namingMap = map[string]string{
	"Chat": "chat",
}

func naming(s string) string {
	n, ok := namingMap[s]
	if ok {
		return n
	}
	return s
}

func getDBExcludeMask(mask *core_user.FieldMask) []string {
	masks := cloneDefaultExcludeMask()
	if mask != nil {
		for _, m := range mask.Fields {
			_, ok := masks[m]
			if ok {
				delete(masks, m)
			}
		}
	}

	dbMask := make([]string, 0)
	for m, _ := range masks {
		dbMask = append(dbMask, naming(m))
	}
	return dbMask
}

func cloneDefaultExcludeMask() map[string]struct{} {
	m := make(map[string]struct{})
	for k, _ := range defaultExcludeMask {
		m[k] = struct{}{}
	}
	return m
}

var defaultExcludeMask = map[string]struct{}{
	"Chat": {},
}
