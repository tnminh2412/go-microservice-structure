package redis_sender

import (
	"context"
	"encoding/json"
	"go-pkg/helper/context_helper"
	"go-pkg/log/base"
	"go-pkg/redis"

	"github.com/tnminh2412/go-microservice-structure/service/pkg/notify/sender"

	gored "github.com/go-redis/redis/v8"
)

type Service struct {
	redisClient redis.Client
}

func New() *Service {
	redisClient := redis.GetUniversal("Notify")
	return &Service{redisClient: redisClient}
}

func (s *Service) Send(ctx context.Context, params *sender.SendParams) error {
	pubMsg := pubNotifyMessage{
		ClientID: params.ClientID,
		DirectIDs: func() map[string]struct{} {
			ids := make(map[string]struct{})
			for _, i := range params.IDs {
				ids[i] = struct{}{}
			}
			return ids
		}(),
		Data: params.Data,
	}

	value, err := json.Marshal(pubMsg)
	if err != nil {
		return err
	}
	topic := redisTopic.NotifyDirectMessage(params.ClientID)
	err = s.redisClient.XAdd(ctx, &gored.XAddArgs{
		Stream: topic,
		Values: map[string]interface{}{
			"msg": value,
		},
		MaxLen: notifyMessageMaxLen,
		Approx: true,
	}).Err()
	if err != nil {
		context_helper.Logger(ctx).Error(func() string {
			return "redis publish err"
		}, base.Error(err))
	}

	return nil
}

type BroadcastParams struct {
	ClientID string
	Topics   []string
	Data     map[string]string
}

func (s *Service) Broadcast(ctx context.Context, params *sender.BroadcastParams) error {
	pubMsg := pubNotifyMessage{
		ClientID: params.ClientID,
		Topics: func() map[string]struct{} {
			topics := make(map[string]struct{})
			for _, i := range params.Topics {
				topics[i] = struct{}{}
			}
			return topics
		}(),
		Data: params.Data,
	}

	value, err := json.Marshal(pubMsg)
	if err != nil {
		return err
	}
	topic := redisTopic.NotifyMessage(params.ClientID)
	err = s.redisClient.XAdd(ctx, &gored.XAddArgs{
		Stream: topic,
		Values: map[string]interface{}{
			"msg": value,
		},
		MaxLen: notifyMessageMaxLen,
		Approx: true,
	}).Err()
	if err != nil {
		context_helper.Logger(ctx).Error(func() string {
			return "redis publish err"
		}, base.Error(err))
	}

	return nil
}
