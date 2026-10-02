package redis_sender

import (
	"context"
	"encoding/json"
	"go-pkg/helper/context_helper"
	"go-pkg/log/base"

	"github.com/tnminh2412/go-microservice-structure/service/pkg/notify/sender"

	gored "github.com/go-redis/redis/v8"
)

func (s *Service) Register(ctx context.Context, params *sender.RegisterParams) error {

	pubMsg := pubNotifyManagement{
		ClientID:     params.ClientID,
		Event:        event.Register,
		NotifyID:     params.ID,
		Username:     params.Username,
		UserType:     params.UserType,
		CurrencyCode: params.CurrencyCode,
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

func (s *Service) Deregister(ctx context.Context, params *sender.DeregisterParams) error {

	pubMsg := pubNotifyManagement{
		ClientID: params.ClientID,
		Event:    event.DeRegister,
		NotifyID: params.ID,
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
