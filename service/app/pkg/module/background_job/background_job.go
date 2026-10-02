package background_job

import (
	"context"
	"time"
)

type Task struct {
	data interface{}
}

func New(inputData interface{}) *Task {
	return &Task{
		data: inputData,
	}
}

func (t *Task) Run(ctx context.Context, fn func(ctx context.Context)) {
	go func() {
		bgCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
		defer cancel()

		defer func() {
			if r := recover(); r != nil {
			}
		}()

		fn(bgCtx)
	}()
}
