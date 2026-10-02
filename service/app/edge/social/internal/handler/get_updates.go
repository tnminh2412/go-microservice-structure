package handler

import (
	"context"
	"go-pkg/helper/context_helper"
	service_error "go-pkg/helper/error_helper"

	_struct "github.com/golang/protobuf/ptypes/struct"

	"github.com/tnminh2412/go-microservice-structure/service/app/edge/social/internal/pkg/error_helper"
	"github.com/tnminh2412/go-microservice-structure/service/app/edge/social/internal/pkg/request_context"
)

func (s Service) WebbHook(ctx context.Context, req *_struct.Struct) (*_struct.Struct, error) {
	domain := context_helper.Request(ctx).ForwardedHost()
	token := request_context.SecretToken(ctx)

	bytes, err := req.MarshalJSON()
	if err != nil {
		context_helper.Logger(ctx).Error(func() string {
			return "io.ReadAll(r.Body) error"
		})

		return nil, error_helper.New(service_error.Internal, "", nil)
	}

	err = s.botManager.ProcessUpdates(domain, token, bytes)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
