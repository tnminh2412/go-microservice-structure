package error_helper

import (
	"go-pkg/helper/error_helper"
	"proto-repo/generated/cmm/service/cmm/integration/user/integration_user"

	"github.com/cloudflare/cloudflare-go"
)

func New(code string, message string, internal error) *error_helper.ServiceError {
	return error_helper.New(code, error_helper.DomainAPI, Domain, message, internal, "")
}

const Domain = "UserService"

var errorMapping = map[string]string{
	integration_user.ErrorMessage_InvalidPassword.String():           error_helper.InvalidArgument,
	integration_user.ErrorMessage_InvalidUsernameOrPassword.String(): error_helper.InvalidArgument,
	integration_user.ErrorMessage_InvalidStartOrEndTime.String():     error_helper.InvalidArgument,
	integration_user.ErrorMessage_InvalidOtp.String():                error_helper.InvalidArgument,
	integration_user.ErrorMessage_InvalidOtpIsRequired.String():      error_helper.InvalidArgument,
	integration_user.ErrorMessage_InvalidEmail.String():              error_helper.InvalidArgument,
	integration_user.ErrorMessage_InvalidToken.String():              error_helper.InvalidArgument,
	integration_user.ErrorMessage_InvalidTokenExpired.String():       error_helper.InvalidArgument,
	integration_user.ErrorMessage_InvalidTokenCompleted.String():     error_helper.InvalidArgument,
}

func init() {
	error_helper.SetErrorMessageMapping(errorMapping)
}

func IsCloudflareDuplicateErr(err error) bool {
	cfApiErr, ok := err.(*cloudflare.RequestError)
	if ok {
		for _, apiErr := range cfApiErr.Errors() {
			if apiErr.Code == 81057 {
				return true
			}
		}
	}
	return false
}

func IsInvalidDomainErr(err error) bool {
	cfApiErr, ok := err.(*cloudflare.RequestError)
	if ok {
		for _, apiErr := range cfApiErr.Errors() {
			if apiErr.Code == 1004 {
				return true
			}
		}
	}
	return false
}
