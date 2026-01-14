package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/whoaa512/asana-cli/internal/config"
)

func TestGetTeam(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/teams/12345" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("unexpected method: %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"gid":  "12345",
				"name": "Engineering Team",
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

	team, err := client.GetTeam(context.Background(), "12345")
	if err != nil {
		t.Fatalf("GetTeam() error = %v", err)
	}

	if team.GID != "12345" {
		t.Errorf("team.GID = %q, want %q", team.GID, "12345")
	}
	if team.Name != "Engineering Team" {
		t.Errorf("team.Name = %q, want %q", team.Name, "Engineering Team")
	}
}
