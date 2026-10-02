package currency_exchange

import (
	"go-pkg/helper/error_helper"
)

const (
	InvalidCurrency = "InvalidCurrency"
)

func Error(code string, message string, internal error) *error_helper.ServiceError {
	return error_helper.New(code, error_helper.DomainAPI, Domain, message, internal, "")
}

const Domain = "CurrencyExchangePkg"
