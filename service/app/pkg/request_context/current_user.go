package request_context

import "context"

type userKey struct{}

func (r Request) CurrentUsername() string {
	username, ok := r.ctx.Value(userKey{}).(string)
	if ok {
		return username
	}
	return ""
}

func SetCurrentUser(ctx context.Context, username string) context.Context {
	return context.WithValue(ctx, userKey{}, username)
}
