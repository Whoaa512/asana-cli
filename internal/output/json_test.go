package output

import (
	"bytes"
	"encoding/json"
	stderrors "errors"
	"testing"

	"github.com/whoaa512/asana-cli/internal/errors"
)

func TestPrint(t *testing.T) {
	var buf bytes.Buffer
	out := NewJSON(&buf)

	data := map[string]any{
		"gid":  "123",
		"name": "Test Task",
	}

	if err := out.Print(data); err != nil {
		t.Fatalf("Print() error = %v", err)
	}

	var result map[string]any
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}

	if result["gid"] != "123" {
		t.Errorf("gid = %v, want %v", result["gid"], "123")
	}
}

func TestPrintError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		wantCode string
		wantExit int
	}{
		{
			name:     "cli error",
			err:      errors.NewAuthError("invalid token"),
			wantCode: "AUTH_FAILURE",
			wantExit: errors.ExitAuthFailure,
		},
		{
			name:     "not found",
			err:      errors.NewNotFoundError("task"),
			wantCode: "NOT_FOUND",
			wantExit: errors.ExitNotFound,
		},
		{
			name:     "plain error",
			err:      stderrors.New("something went wrong"),
			wantCode: "GENERAL_ERROR",
			wantExit: errors.ExitGeneral,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			out := NewJSON(&buf)

			if err := out.PrintError(tt.err); err != nil {
				t.Fatalf("PrintError() error = %v", err)
			}

			var result ErrorResponse
			if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
				t.Fatalf("invalid JSON: %v", err)
			}

			if result.Error.Code != tt.wantCode {
				t.Errorf("code = %q, want %q", result.Error.Code, tt.wantCode)
			}
			if result.Error.ExitCode != tt.wantExit {
				t.Errorf("exit_code = %d, want %d", result.Error.ExitCode, tt.wantExit)
			}
		})
	}
}
