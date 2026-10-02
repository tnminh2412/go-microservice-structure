package rate_limiter

import (
	"context"
	"go-pkg/redis"
	"time"

	"github.com/go-redis/redis_rate/v9"
)

type Limit struct {
	Rate   int
	Period time.Duration
}
type IService interface {
	Limit(ctx context.Context, key string, limit Limit) (time.Duration, error)
}

type rateLimiter struct {
	limiter *redis_rate.Limiter
}

func New(cluster redis.Client) *rateLimiter {
	limiter := redis_rate.NewLimiter(cluster)
	return &rateLimiter{limiter: limiter}
}
