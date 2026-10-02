package activity

import "context"

type Log interface {
	Append(ctx context.Context, clientID string, activity Activity) error
}
