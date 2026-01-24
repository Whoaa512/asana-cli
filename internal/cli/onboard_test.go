package cli

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/whoaa512/asana-cli/internal/api"
	"github.com/whoaa512/asana-cli/internal/config"
)

func TestOnboardCommandRegistered(t *testing.T) {
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Use == "onboard" {
			found = true
			break
		}
	}
	if !found {
		t.Error("onboard command should be registered")
	}
}

func TestOnboardCommandShort(t *testing.T) {
	if onboardCmd.Short == "" {
		t.Error("onboard command should have short description")
	}
}

func TestValidateToken(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		response   string
		wantErr    bool
	}{
		{
			name:       "valid token",
			statusCode: http.StatusOK,
			response:   `{"data": {"gid": "123", "name": "Test User", "email": "test@example.com"}}`,
			wantErr:    false,
		},
		{
			name:       "invalid token",
			statusCode: http.StatusUnauthorized,
			response:   `{"errors": [{"message": "Not Authorized"}]}`,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				_, _ = w.Write([]byte(tt.response))
			}))
			defer server.Close()

			cfg := &config.Config{
				AccessToken: "test-token",
				Timeout:     config.DefaultTimeout,
			}
			client := api.NewHTTPClient(cfg, api.WithBaseURL(server.URL))

			_, err := client.GetMe(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMe() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPickWorkspace(t *testing.T) {
	tests := []struct {
		name        string
		response    string
		wantErr     bool
		wantAutoGID string
	}{
		{
			name:        "single workspace auto-selects",
			response:    `{"data": [{"gid": "123", "name": "My Workspace", "is_organization": false}]}`,
			wantErr:     false,
			wantAutoGID: "123",
		},
		{
			name:     "no workspaces returns error",
			response: `{"data": []}`,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(tt.response))
			}))
			defer server.Close()

			cfg := &config.Config{
				AccessToken: "test-token",
				Timeout:     config.DefaultTimeout,
			}
			client := api.NewHTTPClient(cfg, api.WithBaseURL(server.URL))

			ws, err := pickWorkspace(context.Background(), client)
			if (err != nil) != tt.wantErr {
				t.Errorf("pickWorkspace() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantAutoGID != "" && ws.GID != tt.wantAutoGID {
				t.Errorf("pickWorkspace() GID = %v, want %v", ws.GID, tt.wantAutoGID)
			}
		})
	}
}

func TestOnboardConfigJSON(t *testing.T) {
	cfg := onboardConfig{
		DefaultWorkspace: "ws-123",
		DefaultTeam:      "team-456",
		DefaultProject:   "proj-789",
	}

	if cfg.DefaultWorkspace != "ws-123" {
		t.Errorf("DefaultWorkspace = %v, want ws-123", cfg.DefaultWorkspace)
	}
	if cfg.DefaultTeam != "team-456" {
		t.Errorf("DefaultTeam = %v, want team-456", cfg.DefaultTeam)
	}
	if cfg.DefaultProject != "proj-789" {
		t.Errorf("DefaultProject = %v, want proj-789", cfg.DefaultProject)
	}
}
