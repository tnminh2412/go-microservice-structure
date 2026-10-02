package cache

import (
	"fmt"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/shopspring/decimal"
)

type Service struct {
	c *cache.Cache
}

func New() *Service {
	c := cache.New(5*time.Minute, 10*time.Minute)
	return &Service{c: c}
}

func rateKey(clientID string) string {
	return fmt.Sprintf("_%s:rate", clientID)
}

func (s *Service) SetRates(clientID string, rates map[string]decimal.Decimal, duration time.Duration) {
	s.c.Set(rateKey(clientID), rates, duration)
}

func (s *Service) GetRates(clientID string) (map[string]decimal.Decimal, bool) {
	v, ok := s.c.Get(rateKey(clientID))
	if !ok {
		return nil, false
	}
	return v.(map[string]decimal.Decimal), true
}
