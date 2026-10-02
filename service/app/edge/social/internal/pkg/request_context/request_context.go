package request_context

import (
	"context"
	"go-pkg/helper/context_helper"
)

const (
	secretTokenHeaderKey = "X-Telegram-Bot-Api-Secret-Token"
)

func SecretToken(ctx context.Context) string {
	vals := context_helper.Request(ctx).InHeader(secretTokenHeaderKey)
	if vals == nil || len(vals) == 0 {
		return ""
	}
	return vals[0]
}
