package notify

import (
	"context"

	"github.com/tnminh2412/go-microservice-structure/service/pkg/notify/sender"
	"github.com/tnminh2412/go-microservice-structure/service/pkg/notify/sender/redis_sender"
)

type IService interface {
	Send(ctx context.Context, params *SendParams) error
	Broadcast(ctx context.Context, params *BroadcastParams) error
	Register(ctx context.Context, params *RegisterParams) error
	Deregister(ctx context.Context, params *DeregisterParams) error
	Subscribe(ctx context.Context, params *SubscribeParams) error
	UnSubscribe(ctx context.Context, params *UnSubscribeParams) error
}

type Service struct {
	sender sender.ISender
}

func New() IService {
	return &Service{sender: redis_sender.New()}
}
