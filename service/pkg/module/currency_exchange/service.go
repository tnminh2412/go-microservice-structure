package currency_exchange

import (
	"context"

	"github.com/valyala/fasthttp"
	"golang.org/x/sync/singleflight"

	"github.com/tnminh2412/go-microservice-structure/service/pkg/module/currency_exchange/cache"
)

type IService interface {
	Exchange(ctx context.Context, req *ExchangeRequest) (*ExchangeResult, error)
}

type Service struct {
	cache    *cache.Service
	fastHttp *fasthttp.Client
	group    singleflight.Group
}

func New() IService {
	return &Service{
		cache:    cache.New(),
		fastHttp: &fasthttp.Client{},
	}
}
