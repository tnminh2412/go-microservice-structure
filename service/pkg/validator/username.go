package validator

import (
	"context"
	service_error "go-pkg/helper/error_helper"
	"regexp"
	"strings"
)

func (v Validator) Username(ctx context.Context, username string, language string) error {
	if len(username) == 0 {
		msg := "Username is required"
		return Error(InvalidUsername, msg, nil)
	}

	rules, err := usernameRules(v.UsernameRules)
	if err != nil {
		return err
	}

	for _, rule := range rules {
		if !rule.Regex.MatchString(username) {
			desc := getUsernameRuleDescription(rule, language)
			msg := desc
			return Error(InvalidUsername, msg, nil)
		}
	}

	for _, u := range v.UsernameBlackList {
		if strings.Contains(username, u) == true {
			return Error(InvalidUsernameForbidden, "Username is forbidden", nil)
		}
	}
	return nil
}

type usernameRule struct {
	Regex   *regexp.Regexp
	Locales []Locale
}

func usernameRules(usrRules []UsernameRule) ([]usernameRule, error) {
	rules := make([]usernameRule, 0)
	for _, r := range usrRules {
		re, err := regexp.Compile(r.Regex)
		if err != nil {
			msg := "unexpected error happen"
			return nil, Error(service_error.Internal, msg, nil)
		}
		rule := usernameRule{
			Regex:   re,
			Locales: r.Locales,
		}
		rules = append(rules, rule)
	}
	return rules, nil
}

func getUsernameRuleDescription(p usernameRule, lang string) string {
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
