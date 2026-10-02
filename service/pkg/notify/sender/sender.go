package sender

import (
	"context"
	"time"
)

type SendParams struct {
	ClientID string
	IDs      []string
	Data     string
}

type BroadcastParams struct {
	ClientID string
	Topics   []string
	Data     string
}

type RegisterParams struct {
	ClientID     string
	ID           string
	Token        string
	Time         time.Time
	Info         map[string]string
	Username     string
	UserType     string
	CurrencyCode string
}

type DeregisterParams struct {
	ClientID string
	ID       string
	Time     time.Time
}

type SubscribeParams struct {
	ClientID        string
	ID              string
	Time            time.Time
	Topics          []string
	IsReplaceTopics bool
}

type UnSubscribeParams struct {
	ClientID string
	ID       string
	Time     time.Time
	Topics   []string
}

type ISender interface {
	Send(ctx context.Context, params *SendParams) error
	Broadcast(ctx context.Context, params *BroadcastParams) error
	Register(ctx context.Context, params *RegisterParams) error
	Deregister(ctx context.Context, params *DeregisterParams) error
	Subscribe(ctx context.Context, params *SubscribeParams) error
	UnSubscribe(ctx context.Context, params *UnSubscribeParams) error
}
