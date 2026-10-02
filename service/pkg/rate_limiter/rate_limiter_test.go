package rate_limiter

import (
	"context"
	"fmt"
	"go-pkg/redis"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRateLimit(t *testing.T) {
	l := New(redis.GetUniversal("RateLimiter"))
	ctx := context.Background()

	Convey("limit should work", t, func() {
		limit := Limit{Rate: 1, Period: 1 * time.Minute}
		testID := fmt.Sprintf("test_id_%v", time.Now().Nanosecond())
		r, err := l.Limit(ctx, testID, limit)
		So(err, ShouldBeNil)
		So(r, ShouldBeZeroValue)
	})

	Convey("limit should cause exhausted err", t, func() {
		limit := Limit{Rate: 1, Period: 1 * time.Minute}
		testID := fmt.Sprintf("test_id_%v", time.Now().Nanosecond())
		r, err := l.Limit(ctx, testID, limit)
		So(err, ShouldBeNil)
		So(r, ShouldBeZeroValue)

		r, err = l.Limit(ctx, testID, limit)
		So(err, ShouldBeError)
		So(r, ShouldBeGreaterThan, 0)
	})
}

func TestLimit(t *testing.T) {
	l := New(redis.GetUniversal("RateLimiter"))
	ctx := context.Background()

	Convey("limit example manual click", t, func() {
		limit := Limit{Rate: 3, Period: 1 * time.Minute}
		testID := fmt.Sprintf("test_id_%v", "test")
		r, err := l.Limit(ctx, testID, limit)
		if err != nil {
			t.Logf("err %v", err)
		} else {
			t.Log("limit allowed")
		}
		if r > 0 {
			t.Logf("retry after %v", r)
		}
	})

}
