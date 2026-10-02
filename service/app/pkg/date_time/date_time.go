package date_time

import (
	"errors"
	"strconv"
	"time"
)

type DateTime struct {
	datetime int64
	time     time.Time
}

func Parse(date int64, loc *time.Location) (DateTime, error) {
	t, err := time.ParseInLocation("2006010215", strconv.FormatInt(date, 10), loc)
	if err != nil {
		return DateTime{}, err
	}
	return DateTime{datetime: date, time: t}, nil
}

type Options struct {
	AllowRoundingDateTime *bool
}

func (o *Options) SetAllowRoundingDateTime(b bool) *Options {
	o.AllowRoundingDateTime = &b
	return o
}

var no = false

func mergeOptions(opts ...*Options) *Options {
	o := &Options{
		AllowRoundingDateTime: &no,
	}
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if opt.AllowRoundingDateTime != nil {
			o.AllowRoundingDateTime = opt.AllowRoundingDateTime
		}
	}
	return o
}
func FromTime(in time.Time, loc *time.Location, opts ...*Options) (DateTime, error) {

	opt := mergeOptions(opts...)

	exactDateTime := !*opt.AllowRoundingDateTime

	t := in.In(loc)
	d := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, loc)
	if exactDateTime && !t.Equal(d) {
		return DateTime{}, errors.New("invalid date time")
	}
	f := t.Format("2006010215")
	i, err := strconv.Atoi(f)
	if err != nil {
		return DateTime{}, err
	}
	return DateTime{datetime: int64(i), time: d}, nil
}

func (d DateTime) AsDateTime() int64 {
	return d.datetime
}
func (d DateTime) AsTime() time.Time {
	return d.time
}
