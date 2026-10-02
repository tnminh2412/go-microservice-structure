package validator

import (
	service_error "go-pkg/helper/error_helper"

	"github.com/tnminh2412/go-microservice-structure/service/app/integration/user/internal/pkg/error_helper"
)

func (v Validator) Page(skip int64, limit int64) error {
	if skip < 0 {
		msg := "Skip must not be negative"
		return error_helper.New(service_error.InvalidArgument, msg, nil)
	}
	if limit <= 0 {
		msg := "Limit must be positive"
		return error_helper.New(service_error.InvalidArgument, msg, nil)
	}
	return nil
}
