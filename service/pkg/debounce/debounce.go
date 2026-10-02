package debounce

import (
	"sync"
	"time"
)

func After(after time.Duration) func(f func()) {
	d := &debouncer{after: after}

	return func(f func()) {
		d.debounce(f)
	}
}

type debouncer struct {
	sync.Mutex
	after time.Duration
	timer *time.Timer
}

func (d *debouncer) debounce(f func()) {
	d.Lock()
	defer d.Unlock()

	if d.timer != nil {
		d.timer.Stop()
	}
	d.timer = time.AfterFunc(d.after, f)
}
