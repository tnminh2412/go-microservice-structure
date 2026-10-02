package scheduler_lock

import (
	"go-pkg/helper/error_helper"
)

const (
	LockFailure          = "LockFailure"
	UnLockFailure        = "UnLockFailure"
	RepeatTimeNotElapsed = "RepeatTimeNotElapsed"
)

func Error(code string, message string, internal error) *error_helper.ServiceError {
	return error_helper.New(code, error_helper.DomainAPI, Domain, message, internal, "")
}

const Domain = "SchedulerLock"
