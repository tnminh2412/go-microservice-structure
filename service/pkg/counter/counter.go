package counter

import (
	"context"

	gored "github.com/go-redis/redis/v8"
)

type IService interface {
	Next(ctx context.Context, id string) (int64, error)
}

type Service struct {
	rdb gored.Cmdable
}

func New(rdb gored.Cmdable) *Service {
	return &Service{rdb: rdb}
}

func (s *Service) Next(ctx context.Context, id string) (int64, error) {
	c, err := s.rdb.Incr(ctx, id).Result()
	if err != nil {
		return 0, err
	}
	return c, nil
}
