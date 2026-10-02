package request_context

import (
	"context"
	"go-pkg/helper/context_helper"

	"github.com/tnminh2412/go-microservice-structure/service/app/edge/backoffice/internal/pkg/request_context/headers"
)

func ClientID(ctx context.Context) string {
	vals := context_helper.Request(ctx).InHeader(headers.ClientIDKey)
	if vals == nil || len(vals) == 0 {
		return ""
	}
	return vals[0]
}

func IP(ctx context.Context) string {
	return context_helper.Request(ctx).CustomerIP()
}

func Country(ctx context.Context) string {
	return context_helper.Request(ctx).Country()
}

func ClientDomain(ctx context.Context) string {
	return context_helper.Request(ctx).ForwardedHost()
}
