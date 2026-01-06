package errors

import (
	"errors"
	"fmt"
)

const (
	ExitSuccess      = 0
	ExitGeneral      = 1
	ExitInvalidArgs  = 2
	ExitAuthFailure  = 3
	ExitNotFound     = 4
	ExitRateLimited  = 5
	ExitNetworkError = 6
)

type CLIError struct {
	Message  string `json:"message"`
	Code     string `json:"code"`
	ExitCode int    `json:"exit_code"`
	Cause    error  `json:"-"`
}

func (e *CLIError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

func (e *CLIError) Unwrap() error {
	return e.Cause
}

func NewGeneralError(msg string, cause error) *CLIError {
	return &CLIError{
		Message:  msg,
		Code:     "GENERAL_ERROR",
		ExitCode: ExitGeneral,
		Cause:    cause,
	}
}

func NewInvalidArgsError(msg string) *CLIError {
	return &CLIError{
		Message:  msg,
		Code:     "INVALID_ARGS",
		ExitCode: ExitInvalidArgs,
	}
}

func NewAuthError(msg string) *CLIError {
	return &CLIError{
		Message:  msg,
		Code:     "AUTH_FAILURE",
		ExitCode: ExitAuthFailure,
	}
}

func NewNotFoundError(resource string) *CLIError {
	return &CLIError{
		Message:  fmt.Sprintf("%s not found", resource),
		Code:     "NOT_FOUND",
		ExitCode: ExitNotFound,
	}
}

func NewRateLimitedError(retryAfter string) *CLIError {
	msg := "rate limited"
	if retryAfter != "" {
		msg = fmt.Sprintf("rate limited, retry after %s", retryAfter)
	}
	return &CLIError{
		Message:  msg,
		Code:     "RATE_LIMITED",
		ExitCode: ExitRateLimited,
	}
}

func NewNetworkError(msg string, cause error) *CLIError {
	return &CLIError{
		Message:  msg,
		Code:     "NETWORK_ERROR",
		ExitCode: ExitNetworkError,
		Cause:    cause,
	}
}

func GetExitCode(err error) int {
	var cliErr *CLIError
	if errors.As(err, &cliErr) {
		return cliErr.ExitCode
	}
	return ExitGeneral
}

func AsCLIError(err error) *CLIError {
	var cliErr *CLIError
	if errors.As(err, &cliErr) {
		return cliErr
	}
	return NewGeneralError(err.Error(), err)
}
