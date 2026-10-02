package handler

import (
	"context"
	"fmt"
	"go-pkg/helper/context_helper"
	service_error "go-pkg/helper/error_helper"
	"go-pkg/log"
	"go-pkg/log/base"
	"sync"
	"time"

	"github.com/tnminh2412/go-microservice-structure/service/app/integration/pkg/feature/person/player/token_player"
	"github.com/tnminh2412/go-microservice-structure/service/app/integration/worker/player_signout/internal/pkg/error_helper"
	"github.com/tnminh2412/go-microservice-structure/service/app/integration/worker/player_signout/internal/pkg/local_conf"
	"github.com/tnminh2412/go-microservice-structure/service/app/pkg/shared_conf"
	"github.com/tnminh2412/go-microservice-structure/service/pkg/global_conf"
	"github.com/tnminh2412/go-microservice-structure/service/pkg/notify"
	"github.com/tnminh2412/go-microservice-structure/service/pkg/scheduler_lock"
)

type playerSignoutTask struct {
	ctx         context.Context
	cancel      context.CancelFunc
	repeatAfter time.Duration
	startAt     time.Time
	tokenSrv    *token_player.Server
	globalConf  global_conf.IService
	localConf   local_conf.IService
	sharedConf  shared_conf.IService
	locker      scheduler_lock.IService
	notifier    notify.IService
}

func (t playerSignoutTask) Context() context.Context     { return t.ctx }
func (t playerSignoutTask) GetRepeatTime() time.Duration { return t.repeatAfter }
func (t playerSignoutTask) Stop() <-chan struct{}        { return t.ctx.Done() }
func (t playerSignoutTask) StartAt() time.Time           { return t.startAt }
func (t playerSignoutTask) Type() string                 { return "PlayerSignout" }
func (t playerSignoutTask) Name() string                 { return "Worker:PlayerSignout" }
func (t playerSignoutTask) Clean()                       {}
func (t playerSignoutTask) CancelFunc() context.CancelFunc {
	_, cancelFunc := context.WithCancel(t.ctx)
	return cancelFunc
}
func (t playerSignoutTask) KillIfErr() bool        { return false }
func (t playerSignoutTask) MaxRetry() int64        { return -1 }
func (t playerSignoutTask) RetryIn() time.Duration { return time.Second * 5 }
func (t playerSignoutTask) Expire() time.Duration  { return 5 * time.Minute }
func (t playerSignoutTask) ErrorHandler(err error) {
	log.Error(func() string {
		return fmt.Sprintf("task %s got err: %v", t.Name(), err)
	})
}
func (t playerSignoutTask) TimeToFirstRun() time.Duration { return 5 * time.Second }
func (t playerSignoutTask) DailyMaxRun() int              { return 0 }

func (t playerSignoutTask) Run(ctx context.Context) error {

	context_helper.Logger(ctx).Info(func() string {
		return "process start"
	})

	globalConfigs, err := t.globalConf.GetAllClientConf(ctx)
	if err != nil {
		context_helper.Logger(ctx).Error(func() string {
			return "get conf err"
		}, base.Error(err))
		return err
	}
	now := time.Now()

	var wg sync.WaitGroup
	errCh := make(chan error, len(globalConfigs))

	for clientID, c := range globalConfigs {
		wg.Add(1)

		go func(clientID string, c *global_conf.Conf) {

			defer wg.Done()

			op := fmt.Sprintf("_%v:task:player_signout", clientID)
			context_helper.Logger(ctx).Debug(func() string {
				return "PlayerSignOut"
			}, base.Field("clientID", clientID))
			err := t.process(ctx, clientID, op, now)
			if err != nil {
				context_helper.Logger(ctx).Error(func() string {
					return "t.process error"
				}, base.Error(err), base.Field("clientID", clientID))
				errCh <- error_helper.New(service_error.Internal, "", err)
			}
		}(clientID, c)
	}

	wg.Wait()

	context_helper.Logger(ctx).Info(func() string {
		return "process finish"
	})

	close(errCh)
	for err := range errCh {
		context_helper.Logger(ctx).Error(func() string {
			return "clear inactive player err"
		}, base.Error(err))
	}
	return nil

}

func (t playerSignoutTask) process(ctx context.Context, clientID string, taskName string, now time.Time) error {

	locConf, err := t.localConf.GetConfigValue(ctx, clientID)
	if err != nil {
		return error_helper.New(service_error.Internal, "", err)
	}

	finish, err := t.locker.Lock(ctx, taskName, time.Duration(locConf.RepeatSeconds)*time.Second)
	if err != nil {
		if _, ok := service_error.Is(err, scheduler_lock.RepeatTimeNotElapsed); ok {
			context_helper.Logger(ctx).Debug(func() string {
				return fmt.Sprintf("run time for %v not elapse", clientID)
			})
			return nil
		}
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, locConf.RepeatSeconds*time.Second)
	defer cancel()

	defer func() {
		err = finish(ctx)
		if err != nil {
			context_helper.Logger(ctx).Error(func() string {
				return "finish fail"
			}, base.Error(err))
		}
	}()

	context_helper.Logger(ctx).Info(func() string {
		return "start inactive player clear"
	}, base.Field("clientID", clientID))

	params := token_player.ClearInactivePlayerParams{
		Limit:         locConf.BatchSize,
		MinActiveTime: now.Add(-locConf.OnlineBackDurationSeconds * time.Second),
	}
	hasMore := true
	count := int64(0)
	for hasMore {
		res, err := t.tokenSrv.ClearInactivePlayer(ctx, clientID, params)
		if err != nil {
			return error_helper.New(service_error.Internal, "", err)
		}

		notifyIDs := make([]string, 0)
		for _, u := range res.Usernames {
			id := notify.ID.Player(clientID, u)
			notifyIDs = append(notifyIDs, id)

			err = t.notifier.Deregister(ctx, &notify.DeregisterParams{
				ClientID: clientID,
				ID:       id,
				Time:     now,
			})
			if err != nil {
				context_helper.Logger(ctx).Error(func() string {
					return "deregister fail"
				}, base.Error(err))
			}
		}
		err = t.notifier.Send(ctx, &notify.SendParams{
			ClientID: clientID,
			IDs:      notifyIDs,
			Data:     notify.SignoutData(""),
		})
		if err != nil {
			context_helper.Logger(ctx).Error(func() string {
				return "notify fail"
			}, base.Error(err))
		}

		hasMore = res.Count == params.Limit
		count += res.Count
	}

	context_helper.Logger(ctx).Info(func() string {
		return fmt.Sprintf("complete inactive player clear - count %v", count)
	}, base.Field("clientID", clientID))
	return nil
}
