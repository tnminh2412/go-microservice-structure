package background_log

import (
	"context"
	"go-pkg/helper/context_helper"
	"go-pkg/log/base"
	"time"

	"github.com/tnminh2412/go-microservice-structure/service/app/edge/backoffice/internal/pkg/activity"
)

type noCancel struct {
	ctx context.Context
}

func (c noCancel) Deadline() (time.Time, bool)       { return time.Time{}, false }
func (c noCancel) Done() <-chan struct{}             { return nil }
func (c noCancel) Err() error                        { return nil }
func (c noCancel) Value(key interface{}) interface{} { return c.ctx.Value(key) }

// withoutCancel returns a context that is never canceled.
func withoutCancel(ctx context.Context) context.Context {
	return noCancel{ctx: ctx}
}

type Log struct {
	activityLog activity.Log
}

func New(activityLog activity.Log) Log {
	return Log{
		activityLog: activityLog,
	}
}

func (l Log) Append(ctx context.Context, clientID string, activity activity.Activity) error {
	go func() {
		ctx := withoutCancel(ctx)
		ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
		defer cancel()
		err := l.activityLog.Append(ctx, clientID, activity)
		if err != nil {
			context_helper.Logger(ctx).Error(func() string {
				return "fail log activity"
			}, base.Error(err), base.Field("log", activity))
		}
	}()
	return nil
}
