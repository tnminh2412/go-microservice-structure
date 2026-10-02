package date

import (
	"errors"
	"strconv"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Date struct {
	date int64
	time time.Time
}

func Parse(date int64, loc *time.Location) (Date, error) {
	t, err := time.ParseInLocation("20060102", strconv.FormatInt(date, 10), loc)
	if err != nil {
		return Date{}, err
	}
	return Date{date: date, time: t}, nil
}

type Options struct {
	AllowRoundingDate *bool
}

func (o *Options) SetAllowRoundingDate(b bool) *Options {
	o.AllowRoundingDate = &b
	return o
}

var no = false

func mergeOptions(opts ...*Options) *Options {
	o := &Options{
		AllowRoundingDate: &no,
	}
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if opt.AllowRoundingDate != nil {
			o.AllowRoundingDate = opt.AllowRoundingDate
		}
	}
	return o
}
func FromTime(in time.Time, loc *time.Location, opts ...*Options) (Date, error) {

	opt := mergeOptions(opts...)

	exactDate := !*opt.AllowRoundingDate

	t := in.In(loc)
	d := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, loc)
	if exactDate && !t.Equal(d) {
		return Date{}, errors.New("invalid date")
	}
	f := t.Format("20060102")
	i, err := strconv.Atoi(f)
	if err != nil {
		return Date{}, err
	}
	return Date{date: int64(i), time: d}, nil
}

func (d Date) AsDate() int64 {
	return d.date
}
func (d Date) AsTime() time.Time {
	return d.time
}
func (d Date) AsTimestamp() *timestamp.Timestamp {
	return timestamppb.New(d.time)
}
func (d Date) After(t Date) bool {
	return d.AsTime().After(t.AsTime())
}

func (d Date) AddDate(years int, months int, days int) Date {
	t := d.AsTime().AddDate(years, months, days)
	n, _ := FromTime(t, t.Location())
	return n
}
