package time_function

import "time"

var dayOffsets = map[time.Weekday]int{
	time.Monday:    0,
	time.Tuesday:   1,
	time.Wednesday: 2,
	time.Thursday:  3,
	time.Friday:    4,
	time.Saturday:  5,
	time.Sunday:    6,
}

func BeginningOfWeek(t time.Time) time.Time {
	dayOffset := dayOffsets[t.Weekday()]
	targetDate := t.AddDate(0, 0, -int(dayOffset))
	return time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(), 0, 0, 0, 0, t.Location())
}
