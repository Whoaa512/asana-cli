package errors

import (
	"errors"
	"testing"
)

func TestCLIErrorError(t *testing.T) {
	tests := []struct {
		name string
		err  *CLIError
		want string
	}{
		{
			name: "without cause",
			err:  NewAuthError("invalid token"),
			want: "invalid token",
		},
		{
			name: "with cause",
			err:  NewNetworkError("connection failed", errors.New("timeout")),
			want: "connection failed: timeout",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("Error() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestErrorFactories(t *testing.T) {
	tests := []struct {
		name     string
		err      *CLIError
		wantCode string
		wantExit int
	}{
		{"general", NewGeneralError("oops", nil), "GENERAL_ERROR", ExitGeneral},
		{"invalid args", NewInvalidArgsError("bad flag"), "INVALID_ARGS", ExitInvalidArgs},
		{"auth", NewAuthError("bad token"), "AUTH_FAILURE", ExitAuthFailure},
		{"not found", NewNotFoundError("task"), "NOT_FOUND", ExitNotFound},
		{"rate limited", NewRateLimitedError("60s"), "RATE_LIMITED", ExitRateLimited},
		{"network", NewNetworkError("timeout", nil), "NETWORK_ERROR", ExitNetworkError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Code != tt.wantCode {
				t.Errorf("Code = %q, want %q", tt.err.Code, tt.wantCode)
			}
			if tt.err.ExitCode != tt.wantExit {
				t.Errorf("ExitCode = %d, want %d", tt.err.ExitCode, tt.wantExit)
			}
		})
	}
}

func TestGetExitCode(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want int
	}{
		{"cli error", NewAuthError("bad"), ExitAuthFailure},
		{"wrapped cli error", errors.Join(errors.New("wrap"), NewNotFoundError("x")), ExitNotFound},
		{"plain error", errors.New("plain"), ExitGeneral},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetExitCode(tt.err); got != tt.want {
				t.Errorf("GetExitCode() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestAsCLIError(t *testing.T) {
	cliErr := NewAuthError("test")
	plain := errors.New("plain error")

	if got := AsCLIError(cliErr); got != cliErr {
		t.Error("AsCLIError should return same CLIError")
	}

	converted := AsCLIError(plain)
	if converted.Code != "GENERAL_ERROR" {
		t.Errorf("converted Code = %q, want GENERAL_ERROR", converted.Code)
	}
	if converted.Message != "plain error" {
		t.Errorf("converted Message = %q, want 'plain error'", converted.Message)
	}
}
