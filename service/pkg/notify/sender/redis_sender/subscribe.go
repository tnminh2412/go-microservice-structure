package redis_sender

import (
	"context"
	"encoding/json"
	"go-pkg/helper/context_helper"
	"go-pkg/helper/string_helper"
	"go-pkg/log/base"

	"github.com/tnminh2412/go-microservice-structure/service/pkg/notify/sender"

	gored "github.com/go-redis/redis/v8"
)

func (s *Service) Subscribe(ctx context.Context, params *sender.SubscribeParams) error {

	isReplaceTopics := "true"
	if !params.IsReplaceTopics {
		isReplaceTopics = "false"
	}
	pubMsg := pubNotifyManagement{
		ClientID: params.ClientID,
		Event:    event.Subscribe,
		NotifyID: params.ID,
		Data: map[string]string{
			"topics":          string_helper.Combine(params.Topics),
			"isReplaceTopics": isReplaceTopics,
		},
	}

	value, err := json.Marshal(pubMsg)
	if err != nil {
		return err
	}
	topic := redisTopic.NotifyManagement(params.ClientID)
	err = s.redisClient.XAdd(ctx, &gored.XAddArgs{
		Stream: topic,
		Values: map[string]interface{}{
			"msg": value,
		},
		MaxLen: notifyManagementMaxLen,
		Approx: true,
	}).Err()
	if err != nil {
		context_helper.Logger(ctx).Error(func() string {
			return "redis publish err"
		}, base.Error(err))
	}

	return nil
}

func (s *Service) UnSubscribe(ctx context.Context, params *sender.UnSubscribeParams) error {

	pubMsg := pubNotifyManagement{
		ClientID: params.ClientID,
		Event:    event.UnSubscribe,
		NotifyID: params.ID,
		Data: map[string]string{
			"topics": string_helper.Combine(params.Topics),
		},
	}

	value, err := json.Marshal(pubMsg)
	if err != nil {
		return err
	}
	topic := redisTopic.NotifyManagement(params.ClientID)
	err = s.redisClient.XAdd(ctx, &gored.XAddArgs{
		Stream: topic,
		Values: map[string]interface{}{
			"msg": value,
		},
		MaxLen: notifyManagementMaxLen,
		Approx: true,
	}).Err()
	if err != nil {
		context_helper.Logger(ctx).Error(func() string {
			return "redis publish err"
		}, base.Error(err))
	}

	return nil
}
