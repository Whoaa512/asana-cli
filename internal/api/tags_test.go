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

func TestCreateTag(t *testing.T) {
	tests := []struct {
		name    string
		req     TagCreateRequest
		wantErr bool
	}{
		{
			name: "create tag with name only",
			req: TagCreateRequest{
				Name:      "urgent",
				Workspace: "12345",
			},
			wantErr: false,
		},
		{
			name: "create tag with name and color",
			req: TagCreateRequest{
				Name:      "bug",
				Workspace: "12345",
				Color:     "dark-red",
			},
			wantErr: false,
		},
		{
			name: "missing workspace",
			req: TagCreateRequest{
				Name: "test",
			},
			wantErr: true,
		},
		{
			name: "missing name",
			req: TagCreateRequest{
				Workspace: "12345",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				cfg := &config.Config{
					AccessToken: "test-token",
					Timeout:     5 * time.Second,
				}
				client := NewHTTPClient(cfg)
				_, err := client.CreateTag(context.Background(), tt.req)
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				expectedPath := "/workspaces/" + tt.req.Workspace + "/tags"
				if r.URL.Path != expectedPath {
					t.Errorf("unexpected path: %s, want %s", r.URL.Path, expectedPath)
				}
				if r.Method != http.MethodPost {
					t.Errorf("unexpected method: %s", r.Method)
				}

				var req struct {
					Data struct {
						Name  string `json:"name"`
						Color string `json:"color,omitempty"`
					} `json:"data"`
				}
				if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
					t.Fatalf("failed to decode request: %v", err)
				}

				if req.Data.Name != tt.req.Name {
					t.Errorf("name = %q, want %q", req.Data.Name, tt.req.Name)
				}
				if tt.req.Color != "" && req.Data.Color != tt.req.Color {
					t.Errorf("color = %q, want %q", req.Data.Color, tt.req.Color)
				}

				w.Header().Set("Content-Type", "application/json")
				response := map[string]any{
					"data": map[string]any{
						"gid":  "tag123",
						"name": tt.req.Name,
					},
				}
				if tt.req.Color != "" {
					response["data"].(map[string]any)["color"] = tt.req.Color
				}
				err := json.NewEncoder(w).Encode(response)
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

			tag, err := client.CreateTag(context.Background(), tt.req)
			if err != nil {
				t.Fatalf("CreateTag() error = %v", err)
			}

			if tag.GID != "tag123" {
				t.Errorf("tag.GID = %q, want %q", tag.GID, "tag123")
			}
			if tag.Name != tt.req.Name {
				t.Errorf("tag.Name = %q, want %q", tag.Name, tt.req.Name)
			}
			if tt.req.Color != "" && tag.Color != tt.req.Color {
				t.Errorf("tag.Color = %q, want %q", tag.Color, tt.req.Color)
			}
		})
	}
}
