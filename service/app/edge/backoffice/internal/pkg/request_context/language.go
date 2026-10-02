package request_context

import (
	"context"
	"go-pkg/helper/context_helper"
)

func Language(ctx context.Context) string {
	lang := context_helper.Request(ctx).Language()
	if lang == "" {
		return "en"
	}
	return lang
}
