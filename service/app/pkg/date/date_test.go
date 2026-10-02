package date

import (
	"fmt"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func dateTime(s string) time.Time {
	t, _ := time.Parse(time.RFC3339, s)
	return t
}
func TestDate(t *testing.T) {
	Convey("Return err when date is invalid", t, func() {
		cases := map[string]int64{
			"0":           0,
			"invalid day": 20220229,
		}

		for description, test := range cases {
			Convey(fmt.Sprintf("The '%s' test case should be invalid", description), func() {
				_, err := Parse(test, time.UTC)
				So(err, ShouldBeError)

			})
		}
	})

	Convey("Return no err when date is valid", t, func() {
		cases := map[string]int64{
			"valid date": 20220527,
		}

		for description, test := range cases {
			Convey(fmt.Sprintf("The '%s' test case should be valid", description), func() {
				ti, err := Parse(test, time.UTC)
				So(err, ShouldBeNil)
				t.Log(ti.time.String())
			})
		}
	})
}

func TestFromTime(t *testing.T) {
	Convey("Return err when date is invalid", t, func() {
		cases := map[string]time.Time{
			"not exact day": dateTime("2022-01-01T00:00:01Z"),
		}

		for description, test := range cases {
			Convey(fmt.Sprintf("The '%s' test case should be invalid", description), func() {
				_, err := FromTime(test, time.UTC)
				So(err, ShouldBeError)

			})
		}
	})

	Convey("Return no err when date is valid", t, func() {
		cases := map[string]time.Time{
			"valid date": dateTime("2022-01-01T00:00:00Z"),
		}

		for description, test := range cases {
			Convey(fmt.Sprintf("The '%s' test case should be valid", description), func() {
				ti, err := FromTime(test, time.UTC)
				So(err, ShouldBeNil)
				t.Log(ti.time.String())
			})
		}
	})

	Convey("Return no err when allow rounding date", t, func() {
		cases := map[string]time.Time{
			"sec":  dateTime("2022-01-01T00:00:01Z"),
			"min":  dateTime("2022-01-01T00:01:01Z"),
			"hour": dateTime("2022-01-01T01:01:01Z"),
		}

		opt := &Options{}
		opt.SetAllowRoundingDate(true)
		for description, test := range cases {
			Convey(fmt.Sprintf("The '%s' test case should be valid", description), func() {
				ti, err := FromTime(test, time.UTC, opt)
				So(err, ShouldBeNil)
				t.Log(ti.time.String())
			})
		}
	})
}
