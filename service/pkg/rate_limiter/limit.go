package rate_limiter

import (
	"context"
	service_error "go-pkg/helper/error_helper"
	"strconv"
	"time"

	"github.com/go-redis/redis_rate/v9"
)

func (l *rateLimiter) Limit(ctx context.Context, key string, limit Limit) (time.Duration, error) {
	r, err := l.limiter.Allow(ctx, key, redis_rate.Limit{
		Rate:   limit.Rate,
		Burst:  limit.Rate,
		Period: limit.Period,
	})
	if err != nil {
		return 0, Error(service_error.Internal, "limit fail", err)
	}
	if r.RetryAfter >= 0 {
		resetInSec := int(r.RetryAfter.Seconds())
		return r.RetryAfter, Error(service_error.ResourceExhausted, strconv.Itoa(resetInSec), nil)
	}
	return 0, nil
}
