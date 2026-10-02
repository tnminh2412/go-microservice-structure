package handler

import (
	"context"
	"proto-repo/generated/cmm/service/cmm/core/user/core_user"

	"github.com/tnminh2412/go-microservice-structure/service/app/core/user/internal/pkg/local_conf"
	"github.com/tnminh2412/go-microservice-structure/service/app/pkg/shared_conf"
	"github.com/tnminh2412/go-microservice-structure/service/pkg/global_conf"
)

type Service struct {
	// Returns codes.Unimplemented for business RPCs stripped from this sample.
	core_user.UnimplementedUserServiceServer

	GlobalService global_conf.IService
	SharedService shared_conf.IService
	LocalService  local_conf.IService
}

func New() *Service {
	globalService := global_conf.New(context.Background())
	sharedService := shared_conf.New(context.Background(), globalService)
	localService := &local_conf.Service{
		Service: sharedService,
	}

	return &Service{
		GlobalService: globalService,
		SharedService: sharedService,
		LocalService:  localService,
	}
}
