package password

import (
	"go-pkg/cipher/pbkdf2"
	service_error "go-pkg/helper/error_helper"
)

func Compare(value string, correct string) error {
	ok, err := pbkdf2.Check(value, correct)
	if err != nil || !ok {
		return Error(InvalidUsernameOrPassword, "", nil)
	}
	return nil
}

func Hash(value string) (string, error) {
	h, err := pbkdf2.Create(value)
	if err != nil {
		return "", Error(service_error.Internal, "", err)
	}
	return h, nil
}
