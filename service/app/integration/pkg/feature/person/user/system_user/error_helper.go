package system_user

import (
	"go-pkg/helper/error_helper"
)

const (
	InvalidStatusLocked       = "InvalidStatusLocked"
	InvalidUsernameOrPassword = "InvalidUsernameOrPassword"
	InvalidPasswordIsRequired = "InvalidPasswordIsRequired"
	InvalidUsernameDuplicated = "InvalidUsernameDuplicated"
	InvalidCurrentPassword    = "InvalidCurrentPassword"
	InvalidOtpIsRequired      = "InvalidOtpIsRequired"
	InvalidOtp                = "InvalidOtp"
)

func Error(code string, message string, internal error) *error_helper.ServiceError {
	return error_helper.New(code, error_helper.DomainAPI, Domain, message, internal, "")
}

const Domain = "SystemUserPkg"
