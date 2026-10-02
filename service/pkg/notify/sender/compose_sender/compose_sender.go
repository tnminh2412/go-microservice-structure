package compose_sender

import (
	"context"

	"github.com/tnminh2412/go-microservice-structure/service/pkg/notify/sender"
)

type Service struct {
	senders []sender.ISender
}

func New(senders []sender.ISender) *Service {
	return &Service{senders: senders}
}

func (s *Service) Send(ctx context.Context, params *sender.SendParams) error {

	errs := make([]error, 0)
	for _, s := range s.senders {
		err := s.Send(ctx, params)
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) != 0 {
		return errs[0]
	}

	return nil
}

func (s *Service) Broadcast(ctx context.Context, params *sender.BroadcastParams) error {

	errs := make([]error, 0)
	for _, s := range s.senders {
		err := s.Broadcast(ctx, params)
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) != 0 {
		return errs[0]
	}

	return nil
}

func (s *Service) Register(ctx context.Context, params *sender.RegisterParams) error {
	errs := make([]error, 0)
	for _, s := range s.senders {
		err := s.Register(ctx, params)
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) != 0 {
		return errs[0]
	}

	return nil
}
