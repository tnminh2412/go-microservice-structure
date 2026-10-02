package exchange_token

import "go-pkg/helper/error_helper"

func Error(code string, message string, internal error) *error_helper.ServiceError {
	return error_helper.New(code, error_helper.DomainAPI, Domain, message, internal, "")
}

const Domain = "ExchangePkg"
