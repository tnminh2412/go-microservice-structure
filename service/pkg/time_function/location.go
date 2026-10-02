package time_function

import (
	"sync"
	"time"
)

var timeLocStore = sync.Map{}

func LoadLocation(name string) (*time.Location, error) {
	v, ok := timeLocStore.Load(name)
	if ok {
		return v.(*time.Location), nil
	}

	loc, err := time.LoadLocation(name)
	if err != nil {
		return nil, err
	}
	timeLocStore.Store(name, loc)
	return loc, nil
}
