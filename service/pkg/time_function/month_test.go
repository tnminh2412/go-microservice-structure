package time_function

import (
	"fmt"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestBeginningOfMonth(t *testing.T) {

	Convey("Return beginning of month", t, func() {
		cases := map[string]struct {
			input time.Time
			want  time.Time
		}{
			"2": {
				input: time.Date(2023, 2, 20, 0, 0, 0, 0, time.UTC),
				want:  time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC),
			},

			"2 end": {
				input: time.Date(2023, 2, 28, 0, 0, 0, 0, time.UTC),
				want:  time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC),
			},
			"3": {
				input: time.Date(2023, 3, 21, 0, 0, 0, 0, time.UTC),
				want:  time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC),
			},
			"3 end": {
				input: time.Date(2023, 3, 31, 0, 0, 0, 0, time.UTC),
				want:  time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC),
			},
		}

		for description, test := range cases {
			Convey(fmt.Sprintf("The '%s' test case should be valid", description), func() {
				act := BeginningOfMonth(test.input)
				So(act, ShouldEqual, test.want)

				t.Log(act)
			})
		}
	})
}

func TestEndOfMonth(t *testing.T) {

	Convey("Return end of month", t, func() {
		cases := map[string]struct {
			input time.Time
			want  time.Time
		}{
			"2": {
				input: time.Date(2023, 2, 20, 0, 0, 0, 0, time.UTC),
				want:  time.Date(2023, 2, 28, 0, 0, 0, 0, time.UTC),
			},

			"2 end": {
				input: time.Date(2023, 2, 28, 0, 0, 0, 0, time.UTC),
				want:  time.Date(2023, 2, 28, 0, 0, 0, 0, time.UTC),
			},
			"3": {
				input: time.Date(2023, 3, 21, 0, 0, 0, 0, time.UTC),
				want:  time.Date(2023, 3, 31, 0, 0, 0, 0, time.UTC),
			},
			"3 end": {
				input: time.Date(2023, 3, 31, 0, 0, 0, 0, time.UTC),
				want:  time.Date(2023, 3, 31, 0, 0, 0, 0, time.UTC),
			},
		}

		for description, test := range cases {
			Convey(fmt.Sprintf("The '%s' test case should be valid", description), func() {
				act := EndOfMonth(test.input)
				So(act, ShouldEqual, test.want)

				t.Log(act)
			})
		}
	})
}
