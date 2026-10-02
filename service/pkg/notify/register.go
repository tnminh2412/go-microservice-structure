package notify

import (
	"context"
	"time"

	"github.com/tnminh2412/go-microservice-structure/service/pkg/notify/sender"
)

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

func (s *Service) Register(ctx context.Context, params *RegisterParams) error {
	err := s.sender.Register(ctx, &sender.RegisterParams{
		ClientID:     params.ClientID,
		ID:           params.ID,
		Token:        params.Token,
		Time:         params.Time,
		Info:         params.Info,
		Username:     params.Username,
		UserType:     params.UserType,
		CurrencyCode: params.CurrencyCode,
	})
	if err != nil {
		return err
	}
	return nil
}

type DeregisterParams struct {
	ClientID string
	ID       string
	Time     time.Time
}

func (s *Service) Deregister(ctx context.Context, params *DeregisterParams) error {
	err := s.sender.Deregister(ctx, &sender.DeregisterParams{
		ClientID: params.ClientID,
		ID:       params.ID,
		Time:     params.Time,
	})
	if err != nil {
		return err
	}
	return nil
}
