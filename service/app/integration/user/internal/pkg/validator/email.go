package validator

import (
	"net/mail"
	"proto-repo/generated/cmm/service/cmm/integration/user/integration_user"

	"github.com/tnminh2412/go-microservice-structure/service/app/integration/user/internal/pkg/error_helper"
)

func (v Validator) Email(email string) (string, error) {
	addr, err := mail.ParseAddress(email)
	if err != nil {
		msg := "Email is invalid"
		return "", error_helper.New(integration_user.ErrorMessage_InvalidEmail.String(), msg, err)
	}

	return v.emailNormalizer.Normalize(addr.Address), nil
}
