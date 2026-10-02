package log_helper

import (
	"context"
	"fmt"
	"go-pkg/helper/context_helper"
	"go-pkg/log/base"
	"time"
)

func WarnLongElapse(ctx context.Context, warnAfter time.Duration, what string) func() {
	start := time.Now()
	return func() {
		duration := time.Since(start)
		if duration > warnAfter {
			context_helper.Logger(ctx).Warn(func() string {
				return fmt.Sprintf("%s took %v", what, duration)
			}, base.Field("duration", duration))
		}

	}
}
