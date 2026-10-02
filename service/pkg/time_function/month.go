package time_function

import "time"

func BeginningOfMonth(t time.Time) time.Time {
	year, month, _ := t.Date()
	return time.Date(year, month, 1, 0, 0, 0, 0, t.Location())
}

func EndOfMonth(t time.Time) time.Time {
	bom := BeginningOfMonth(t)
	eom := bom.AddDate(0, 1, -1)
	return eom
}
