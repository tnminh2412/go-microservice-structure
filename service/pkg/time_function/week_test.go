package time_function

import (
	"fmt"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestBeginningOfWeek(t *testing.T) {

	Convey("Return beginning of week", t, func() {
		cases := map[string]struct {
			input time.Time
			want  time.Time
		}{
			"mon": {
				input: time.Date(2023, 2, 20, 0, 0, 0, 0, time.UTC),
				want:  time.Date(2023, 2, 20, 0, 0, 0, 0, time.UTC),
			},
			"tues": {
				input: time.Date(2023, 2, 21, 0, 0, 0, 0, time.UTC),
				want:  time.Date(2023, 2, 20, 0, 0, 0, 0, time.UTC),
			},
			"wed": {
				input: time.Date(2023, 2, 22, 0, 0, 0, 0, time.UTC),
				want:  time.Date(2023, 2, 20, 0, 0, 0, 0, time.UTC),
			},
			"thu": {
				input: time.Date(2023, 2, 23, 0, 0, 0, 0, time.UTC),
				want:  time.Date(2023, 2, 20, 0, 0, 0, 0, time.UTC),
			},
			"fri": {
				input: time.Date(2023, 2, 24, 0, 0, 0, 0, time.UTC),
				want:  time.Date(2023, 2, 20, 0, 0, 0, 0, time.UTC),
			},
			"sat": {
				input: time.Date(2023, 2, 25, 0, 0, 0, 0, time.UTC),
				want:  time.Date(2023, 2, 20, 0, 0, 0, 0, time.UTC),
			},
			"sun": {
				input: time.Date(2023, 2, 26, 0, 0, 0, 0, time.UTC),
				want:  time.Date(2023, 2, 20, 0, 0, 0, 0, time.UTC),
			},
			"mon 2": {
				input: time.Date(2023, 2, 27, 0, 0, 0, 0, time.Local),
				want:  time.Date(2023, 2, 27, 0, 0, 0, 0, time.Local),
			},
		}

		for description, test := range cases {
			Convey(fmt.Sprintf("The '%s' test case should be valid", description), func() {
				act := BeginningOfWeek(test.input)
				So(act, ShouldEqual, test.want)

				t.Log(act)
			})
		}
	})
}
