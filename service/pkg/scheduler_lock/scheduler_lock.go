package scheduler_lock

import (
	"context"
	"fmt"
	"strconv"
	"time"

	service_error "go-pkg/helper/error_helper"

	gored "github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
)

type FinishFn func(context.Context) error
type IService interface {
	Lock(ctx context.Context, taskName string, repeatDuration time.Duration) (FinishFn, error)
}

type Service struct {
	rdb gored.UniversalClient
	rs  *redsync.Redsync
}

func New(rdb gored.UniversalClient) *Service {
	pool := goredis.NewPool(rdb)
	rs := redsync.New(pool)

	return &Service{
		rdb: rdb,
		rs:  rs,
	}
}

func lastRunTimeKey(key string) string {
	return fmt.Sprintf("%v:last_run", key)
}
func lockKey(key string) string {
	return fmt.Sprintf("%v:lock", key)
}
func (s *Service) Lock(ctx context.Context, taskName string, repeatDuration time.Duration) (FinishFn, error) {

	v, err := s.rdb.Get(ctx, lastRunTimeKey(taskName)).Result()
	if err != nil {
		if err != gored.Nil {
			return nil, Error(service_error.Database, "", err)
		}
		v = "0"
	}

	i, err := strconv.ParseInt(v, 10, 0)
	if err != nil {
		return nil, Error(service_error.Internal, "parse run time err", err)
	}

	lastRunTime := time.Unix(0, i)
	now := time.Now()
	if now.Add(-repeatDuration).Before(lastRunTime) {
		return nil, Error(RepeatTimeNotElapsed, "", err)
	}

	opts := redsync.WithExpiry(repeatDuration)
	mutex := s.rs.NewMutex(lockKey(taskName), opts)
	err = mutex.LockContext(ctx)
	if err != nil {
		return nil, Error(LockFailure, "", err)
	}

	finishFn := func(ctx context.Context) error {
		ok, err := mutex.UnlockContext(ctx)
		if !ok || err != nil {
			return Error(UnLockFailure, "", err)
		}
		err = s.rdb.Set(ctx, lastRunTimeKey(taskName), now.UnixNano(), 0).Err()
		if err != nil {
			return Error(service_error.Database, "set last run time err", err)
		}
		return nil
	}
	return finishFn, nil
}
