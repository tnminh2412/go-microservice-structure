package scheduler_lock

import (
	"context"
	"go-pkg/redis"
	"testing"
	"time"

	service_error "go-pkg/helper/error_helper"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSchedulerLock(t *testing.T) {
	l := New(redis.GetUniversal("SchedulerLock"))
	ctx := context.Background()

	Convey("lock should work", t, func() {
		finish, err := l.Lock(ctx, "scheduler_test", 1*time.Minute)
		So(err, ShouldBeNil)
		So(finish, ShouldNotBeNil)
		err = finish(ctx)
		So(err, ShouldBeNil)
	})

	Convey("valid run should work", t, func() {
		finish, err := l.Lock(ctx, "scheduler_test_valid_run", 1*time.Minute)
		So(err, ShouldBeNil)
		So(finish, ShouldNotBeNil)
		err = finish(ctx)
		So(err, ShouldBeNil)

		_, err = l.Lock(ctx, "scheduler_test_valid_run", 1*time.Minute)
		So(err, ShouldBeError)
		_, ok := service_error.Is(err, RepeatTimeNotElapsed)
		So(ok, ShouldBeTrue)
	})

}
