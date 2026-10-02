package logger_log

import (
	"context"
	"go-pkg/helper/context_helper"
	"go-pkg/log/base"

	"github.com/tnminh2412/go-microservice-structure/service/app/edge/backoffice/internal/pkg/activity"
)

// Log writes activities to the structured logger. Swap it for a persistent
// sink (e.g. a core activity service) by implementing activity.Log.
type Log struct{}

func New() Log {
	return Log{}
}

func (l Log) Append(ctx context.Context, clientID string, a activity.Activity) error {
	context_helper.Logger(ctx).Info(func() string {
		return "activity"
	},
		base.Field("clientID", clientID),
		base.Field("username", a.Username),
		base.Field("activityID", a.ActivityID),
		base.Field("ip", a.IP),
		base.Field("country", a.Country),
		base.Field("data", a.Data),
	)
	return nil
}
