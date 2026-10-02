package validator

import (
	"net/mail"
)

func (v Validator) Email(email string) (string, error) {
	addr, err := mail.ParseAddress(email)
	if err != nil {
		return "", Error(InvalidEmail, "", err)
	}

	return v.emailNormalizer.Normalize(addr.Address), nil
}
