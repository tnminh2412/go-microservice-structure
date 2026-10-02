package validator

import (
	normalizer "github.com/dimuska139/go-email-normalizer"
)

type Locale struct {
	Language    string
	Description string
}
type PasswordRule struct {
	Regex   string
	Locales []Locale
}

type UsernameRule struct {
	Regex   string
	Locales []Locale
}

type Validator struct {
	UsernameRules     []UsernameRule
	PasswordRules     []PasswordRule
	UsernameBlackList []string
	emailNormalizer   *normalizer.Normalizer
}

func New(passwordRules []PasswordRule, usernameRules []UsernameRule, usernameBlackList []string) Validator {
	return Validator{
		UsernameRules:     usernameRules,
		PasswordRules:     passwordRules,
		UsernameBlackList: usernameBlackList,
		emailNormalizer:   normalizer.NewNormalizer(),
	}
}
