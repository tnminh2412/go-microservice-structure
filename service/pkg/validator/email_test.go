package validator

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestEmail(t *testing.T) {
	Convey("Return err when email is invalid", t, func() {
		cases := map[string]string{
			"empty":             "",
			"no domain":         "abcd",
			"no domain 2":       "abcd@",
			"without id":        "@email.com",
			"invalid character": "test test@email.com",
		}

		v := New([]PasswordRule{}, []UsernameRule{}, []string{})
		for description, test := range cases {
			Convey(fmt.Sprintf("The '%s' test case should be invalid", description), func() {
				_, err := v.Email(test)
				So(err, ShouldBeError)

			})
		}
	})

	Convey("Return normalized email when valid", t, func() {
		cases := map[string]string{
			"domain with underscores": "mail@domain_with_underscores.org",
			"invalid domain":          "test@email",
			"invalid top domain":      "AF_123_def@Aks2.jj2",
			"google domain":           "AF_123_def@gmail.com",
			"google domain 2":         "Af_123_def@gmail.com",
			"google domain 3":         "af_123_def@gmail.com",
			"google domain 4":         "t.e-St+vasya@gmail.com",
		}

		v := New([]PasswordRule{}, []UsernameRule{}, []string{})
		for description, test := range cases {
			Convey(fmt.Sprintf("The '%s' test case should be valid", description), func() {
				addr, err := v.Email(test)
				So(err, ShouldBeNil)

				t.Log(addr)
			})
		}
	})
}
