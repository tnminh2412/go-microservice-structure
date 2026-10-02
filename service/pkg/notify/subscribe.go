package notify

import (
	"context"
	"time"

	"github.com/tnminh2412/go-microservice-structure/service/pkg/notify/sender"
)

type SubscribeParams struct {
	ClientID        string
	ID              string
	Time            time.Time
	Topics          []string
	IsReplaceTopics bool
}

func (s *Service) Subscribe(ctx context.Context, params *SubscribeParams) error {
	err := s.sender.Subscribe(ctx, &sender.SubscribeParams{
		ClientID:        params.ClientID,
		ID:              params.ID,
		Time:            params.Time,
		Topics:          params.Topics,
		IsReplaceTopics: params.IsReplaceTopics,
	})
	if err != nil {
		return err
	}
	return nil
}

type UnSubscribeParams struct {
	ClientID string
	ID       string
	Time     time.Time
	Topics   []string
}

func (s *Service) UnSubscribe(ctx context.Context, params *UnSubscribeParams) error {
	err := s.sender.UnSubscribe(ctx, &sender.UnSubscribeParams{
		ClientID: params.ClientID,
		ID:       params.ID,
		Time:     params.Time,
		Topics:   params.Topics,
	})
	if err != nil {
		return err
	}
	return nil
}
