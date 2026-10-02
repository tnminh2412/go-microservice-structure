package handler

import (
	"context"
	"proto-repo/generated/cmm/service/cmm/integration/worker/player_signout/worker_player_signout"

	"github.com/tnminh2412/go-microservice-structure/service/app/integration/pkg/feature/person/player/token_player"
	"github.com/tnminh2412/go-microservice-structure/service/app/integration/worker/player_signout/internal/global"
	"github.com/tnminh2412/go-microservice-structure/service/app/integration/worker/player_signout/internal/pkg/local_conf"
	"github.com/tnminh2412/go-microservice-structure/service/app/pkg/shared_conf"
	"github.com/tnminh2412/go-microservice-structure/service/pkg/global_conf"
	"github.com/tnminh2412/go-microservice-structure/service/pkg/notify"
	"github.com/tnminh2412/go-microservice-structure/service/pkg/scheduler_lock"
)

type Service struct {
	ctx             context.Context
	tokenSrv        *token_player.Server
	schedulerLocker scheduler_lock.IService
	GlobalService   global_conf.IService
	SharedService   shared_conf.IService
	LocalService    local_conf.IService
	notifier        notify.IService
}

func New() *Service {
	ctx := context.TODO()

	globalService := global_conf.New(ctx)
	sharedService := shared_conf.New(ctx, globalService)
	localService := local_conf.New(sharedService)

	tokenSrv := token_player.New(global.TokenClient())
	schedulerLocker := scheduler_lock.New(global.SchedulerLockClient())
	notifier := notify.New()
	return &Service{
		ctx:             ctx,
		tokenSrv:        tokenSrv,
		schedulerLocker: schedulerLocker,
		GlobalService:   globalService,
		SharedService:   sharedService,
		LocalService:    localService,
		notifier:        notifier,
	}
}

func (s *Service) Ping(ctx context.Context, req *worker_player_signout.PingRequest) (*worker_player_signout.PingReply, error) {
	return &worker_player_signout.PingReply{}, nil
}
