package template_helper

import (
	"encoding/json"
	"time"

	"github.com/shopspring/decimal"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var FuncMap = map[string]any{
	"InTimezone": func(t time.Time, loc string) time.Time {
		l, err := time.LoadLocation(loc)
		if err != nil {
			return t
		}
		return t.In(l)
	},
	"UpperTitle": func(s string) string {
		return cases.Title(language.Und).String(s)
	},
	"ToArrayJson": func(s string) []map[string]interface{} {
		m := make([]map[string]interface{}, 0)
		_ = json.Unmarshal([]byte(s), &m)
		return m
	},
	"ToJson": func(s string) map[string]interface{} {
		m := make(map[string]interface{})
		_ = json.Unmarshal([]byte(s), &m)
		return m
	},
	"ToAmount": func(value string) string {
		amount := decimal.RequireFromString(value)
		str := amount.RoundFloor(8).StringFixed(8)
		return str
	},
}
