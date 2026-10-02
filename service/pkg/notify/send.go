package notify

import (
	"context"

	"github.com/tnminh2412/go-microservice-structure/service/pkg/notify/sender"
)

type SendParams struct {
	ClientID string
	IDs      []string
	Data     string
}

func (s *Service) Send(ctx context.Context, params *SendParams) error {

	err := s.sender.Send(ctx, &sender.SendParams{
		ClientID: params.ClientID,
		IDs:      params.IDs,
		Data:     params.Data,
	})
	if err != nil {
		return err
	}

	return nil
}

type BroadcastParams struct {
	ClientID string
	Topics   []string
	Data     string
}

func (s *Service) Broadcast(ctx context.Context, params *BroadcastParams) error {

	err := s.sender.Broadcast(ctx, &sender.BroadcastParams{
		ClientID: params.ClientID,
		Topics:   params.Topics,
		Data:     params.Data,
	})
	if err != nil {
		return err
	}

	return nil
}
