package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/whoaa512/asana-cli/internal/config"
	"github.com/whoaa512/asana-cli/internal/errors"
)

func TestGetMe(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users/me" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("unexpected method: %s", r.Method)
		}
		if auth := r.Header.Get("Authorization"); auth != "Bearer test-token" {
			t.Errorf("unexpected auth header: %s", auth)
		}

		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"gid":   "12345",
				"name":  "Test User",
				"email": "test@example.com",
			},
		})
		if err != nil {
			t.Fatal(err)
		}
	}))
	defer server.Close()

	cfg := &config.Config{
		AccessToken: "test-token",
		Timeout:     5 * time.Second,
	}
	client := NewHTTPClient(cfg, WithBaseURL(server.URL))

	user, err := client.GetMe(context.Background())
	if err != nil {
		t.Fatalf("GetMe() error = %v", err)
	}

	if user.GID != "12345" {
		t.Errorf("user.GID = %q, want %q", user.GID, "12345")
	}
	if user.Name != "Test User" {
		t.Errorf("user.Name = %q, want %q", user.Name, "Test User")
	}
	if user.Email != "test@example.com" {
		t.Errorf("user.Email = %q, want %q", user.Email, "test@example.com")
	}
}

func TestAPIErrors(t *testing.T) {
	tests := []struct {
		name       string
		status     int
		body       string
		wantExit   int
		wantSubstr string
	}{
		{
			name:       "unauthorized",
			status:     401,
			body:       `{"errors":[{"message":"Not Authorized"}]}`,
			wantExit:   errors.ExitAuthFailure,
			wantSubstr: "Not Authorized",
		},
		{
			name:       "not found",
			status:     404,
			body:       `{"errors":[{"message":"task not found"}]}`,
			wantExit:   errors.ExitNotFound,
			wantSubstr: "not found",
		},
		{
			name:       "rate limited",
			status:     429,
			body:       `{"errors":[{"message":"rate limited"}]}`,
			wantExit:   errors.ExitRateLimited,
			wantSubstr: "rate limited",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(tt.status)
				_, _ = w.Write([]byte(tt.body))
			}))
			defer server.Close()

			cfg := &config.Config{AccessToken: "test", Timeout: 5 * time.Second}
			client := NewHTTPClient(cfg, WithBaseURL(server.URL))

			_, err := client.GetMe(context.Background())
			if err == nil {
				t.Fatal("expected error")
			}

			exitCode := errors.GetExitCode(err)
			if exitCode != tt.wantExit {
				t.Errorf("exit code = %d, want %d", exitCode, tt.wantExit)
			}
			if !strings.Contains(err.Error(), tt.wantSubstr) {
				t.Errorf("error = %q, want to contain %q", err.Error(), tt.wantSubstr)
			}
		})
	}
}

func TestDebugOutput(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"gid":"1","name":"Test"}}`))
	}))
	defer server.Close()

	var debugBuf bytes.Buffer
	cfg := &config.Config{AccessToken: "test-token-12345", Timeout: 5 * time.Second}
	client := NewHTTPClient(cfg, WithBaseURL(server.URL), WithDebug(&debugBuf))

	_, err := client.GetMe(context.Background())
	if err != nil {
		t.Fatalf("GetMe() error = %v", err)
	}

	debugOut := debugBuf.String()
	if !strings.Contains(debugOut, "[DEBUG] GET") {
		t.Error("debug output should contain request method")
	}
	if !strings.Contains(debugOut, "test-tok") {
		t.Error("debug output should contain truncated token")
	}
	if strings.Contains(debugOut, "test-token-12345") {
		t.Error("debug output should NOT contain full token")
	}
}
