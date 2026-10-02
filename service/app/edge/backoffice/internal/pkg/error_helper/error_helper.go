package error_helper

import (
	"encoding/json"
	"go-pkg/helper/error_helper"
	"proto-repo/generated/cmm/service/cmm/edge/backoffice/edge_backoffice"

	"google.golang.org/protobuf/types/known/structpb"
)

func New(code string, message string, internal error) *error_helper.ServiceError {
	return error_helper.New(code, error_helper.DomainAPI, Domain, message, internal, "")
}

type httpErr struct {
	Error *errData
}

type errData struct {
	Code    string
	Message string
	Data    *structpb.Struct
	Child   error
}

func IsHttpServiceError(b []byte) (*error_helper.ServiceError, bool) {
	var e httpErr
	_ = json.Unmarshal(b, &e)
	if e.Error == nil || len(e.Error.Code) == 0 {
		return nil, false
	}
	svcE := e.Error
	return error_helper.New(svcE.Code, Domain, "", svcE.Message, nil, ""), true
}

const Domain = "BackofficeService"

var errorMapping = map[string]string{
	edge_backoffice.ErrorMessage_InvalidUsernameOrPassword.String(): error_helper.InvalidArgument,
	edge_backoffice.ErrorMessage_InvalidPassword.String():           error_helper.InvalidArgument,
	edge_backoffice.ErrorMessage_InvalidCurrentPassword.String():    error_helper.InvalidArgument,
	edge_backoffice.ErrorMessage_InvalidAmount.String():             error_helper.InvalidArgument,
	edge_backoffice.ErrorMessage_InvalidTransactionAction.String():  error_helper.InvalidArgument,
	edge_backoffice.ErrorMessage_MaxLoginAttemptExceed.String():     error_helper.InvalidArgument,
	edge_backoffice.ErrorMessage_InvalidStatus.String():             error_helper.InvalidArgument,
	edge_backoffice.ErrorMessage_InvalidStatusLocked.String():       error_helper.InvalidArgument,
	edge_backoffice.ErrorMessage_InvalidStatusSuspended.String():    error_helper.InvalidArgument,
	edge_backoffice.ErrorMessage_InvalidFileSize.String():           error_helper.InvalidArgument,
	edge_backoffice.ErrorMessage_InvalidFileType.String():           error_helper.InvalidArgument,
	edge_backoffice.ErrorMessage_InvalidFileName.String():           error_helper.InvalidArgument,
	edge_backoffice.ErrorMessage_InvalidStartOrEndTime.String():     error_helper.InvalidArgument,

	edge_backoffice.ErrorMessage_InvalidEmailDuplicated.String():       error_helper.Duplicate,
	edge_backoffice.ErrorMessage_InvalidPhoneNumberDuplicated.String(): error_helper.Duplicate,
	edge_backoffice.ErrorMessage_InvalidUsernameDuplicated.String():    error_helper.Duplicate,

	edge_backoffice.ErrorMessage_InvalidUsername.String():    error_helper.InvalidArgument,
	edge_backoffice.ErrorMessage_InvalidEmail.String():       error_helper.InvalidArgument,
	edge_backoffice.ErrorMessage_InvalidCountryCode.String(): error_helper.InvalidArgument,
	edge_backoffice.ErrorMessage_InvalidPhoneNumber.String(): error_helper.InvalidArgument,

	edge_backoffice.ErrorMessage_InvalidCountryCodeIsRequired.String(): error_helper.InvalidArgument,
	edge_backoffice.ErrorMessage_InvalidEmailIsRequired.String():       error_helper.InvalidArgument,
	edge_backoffice.ErrorMessage_InvalidPhoneNumberIsRequired.String(): error_helper.InvalidArgument,
	edge_backoffice.ErrorMessage_InvalidUsernameIsRequired.String():    error_helper.InvalidArgument,
	edge_backoffice.ErrorMessage_InvalidPasswordIsRequired.String():    error_helper.InvalidArgument,

	edge_backoffice.ErrorMessage_NotEnoughLock.String():       error_helper.InvalidArgument,
	edge_backoffice.ErrorMessage_InvalidLockAmount.String():   error_helper.InvalidArgument,
	edge_backoffice.ErrorMessage_InvalidAmount.String():       error_helper.InvalidArgument,
	edge_backoffice.ErrorMessage_InvalidCurrencyCode.String(): error_helper.InvalidArgument,

	edge_backoffice.ErrorMessage_InvalidOtpIsRequired.String():  error_helper.InvalidArgument,
	edge_backoffice.ErrorMessage_InvalidOtp.String():            error_helper.InvalidArgument,
	edge_backoffice.ErrorMessage_InvalidToken.String():          error_helper.InvalidArgument,
	edge_backoffice.ErrorMessage_InvalidTokenExpired.String():   error_helper.InvalidArgument,
	edge_backoffice.ErrorMessage_InvalidTokenCompleted.String(): error_helper.InvalidArgument,

	edge_backoffice.ErrorMessage_MaxLevelExceed.String():             error_helper.InvalidArgument,
	edge_backoffice.ErrorMessage_PhoneNumberHasBeenVerified.String(): error_helper.InvalidArgument,
	edge_backoffice.ErrorMessage_EmailHasBeenVerified.String():       error_helper.InvalidArgument,
	edge_backoffice.ErrorMessage_InvalidRejoinTimes.String():         error_helper.InvalidArgument,

	edge_backoffice.ErrorMessage_ExceedMinimumAmount.String(): error_helper.InvalidArgument,

	edge_backoffice.ErrorMessage_InvalidAppID.String():     error_helper.InvalidArgument,
	edge_backoffice.ErrorMessage_InvalidHash.String():      error_helper.InvalidArgument,
	edge_backoffice.ErrorMessage_InvalidTimestamp.String(): error_helper.InvalidArgument,
	edge_backoffice.ErrorMessage_InvalidIPAddress.String(): error_helper.InvalidArgument,
}

func init() {
	error_helper.SetErrorMessageMapping(errorMapping)
}
