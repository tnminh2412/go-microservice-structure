package validator

import (
	"context"
	"fmt"
	"go-pkg/helper/context_helper"
	service_error "go-pkg/helper/error_helper"
	"go-pkg/log/base"
	"regexp"
)

const defaultFieldName = "Password"

func (v Validator) Password(ctx context.Context, password string, language string, optionalFieldName ...string) error {
	field := defaultFieldName
	if len(optionalFieldName) > 0 {
		field = optionalFieldName[0]
	}
	if len(password) == 0 {
		msg := fmt.Sprintf("%v is required", field)
		return Error(InvalidPassword, msg, nil)
	}

	rules, err := passwordRules(ctx, v.PasswordRules)
	if err != nil {
		return err
	}

	for _, rule := range rules {
		if !rule.Regex.MatchString(password) {
			desc := getDescription(rule, language)
			msg := desc
			return Error(InvalidPassword, msg, nil)
		}
	}

	return nil
}

type passwordRule struct {
	Regex   *regexp.Regexp
	Locales []Locale
}

func passwordRules(ctx context.Context, pwdRules []PasswordRule) ([]passwordRule, error) {
	rules := make([]passwordRule, 0)
	for _, r := range pwdRules {
		re, err := regexp.Compile(r.Regex)
		if err != nil {
			context_helper.Logger(ctx).Error(func() string {
				return "fail to compile regex"
			}, base.Error(err), base.Field("regex", r.Regex))
			msg := "unexpected error happen"
			return nil, Error(service_error.Internal, msg, nil)
		}
		rule := passwordRule{
			Regex:   re,
			Locales: r.Locales,
		}
		rules = append(rules, rule)
	}
	return rules, nil
}

func getDescription(p passwordRule, lang string) string {
	if len(p.Locales) > 0 {
		for _, l := range p.Locales {
			if l.Language == lang {
				return l.Description
			}
		}
		// by default return first locale
		return p.Locales[0].Description
	}
	return ""
}
