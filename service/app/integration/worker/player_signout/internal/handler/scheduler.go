package handler

import (
	"context"
	"go-pkg/schedule"
	"time"
)

func (s *Service) Start() {
	ctx, cancel := context.WithCancel(context.TODO())
	t := playerSignoutTask{
		ctx:         ctx,
		cancel:      cancel,
		repeatAfter: 1 * time.Minute,
		startAt:     time.Now(),
		tokenSrv:    s.tokenSrv,
		globalConf:  s.GlobalService,
		localConf:   s.LocalService,
		sharedConf:  s.SharedService,
		locker:      s.schedulerLocker,
		notifier:    s.notifier,
	}
	schedule.Register(t, &schedule.Opts{
		ReplaceOld:  true,
		MaxInstance: 1,
	})
}

func (s *Service) Stop() {
	schedule.Stop()
}
