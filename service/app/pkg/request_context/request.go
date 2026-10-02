package request_context

import "context"

type Request struct {
	ctx context.Context
}

func New(ctx context.Context) Request {
	return Request{ctx: ctx}
}
