package error_helper

import (
	"go-pkg/helper/error_helper"
	"proto-repo/generated/cmm/service/cmm/core/user/core_user"
)

func New(code string, message string, internal error) *error_helper.ServiceError {
	return error_helper.New(code, error_helper.DomainAPI, Domain, message, internal, "")
}

const Domain = "UserService"

var grpcErrorMapping = map[string]string{
	core_user.ErrorMessage_InvalidUsernameDuplicated.String():    error_helper.Duplicate,
	core_user.ErrorMessage_InvalidEmailDuplicated.String():       error_helper.Duplicate,
	core_user.ErrorMessage_InvalidPhoneNumberDuplicated.String(): error_helper.Duplicate,
}

func init() {
	error_helper.SetErrorMessageMapping(grpcErrorMapping)
}
