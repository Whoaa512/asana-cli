package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/whoaa512/asana-cli/internal/config"
	"github.com/whoaa512/asana-cli/internal/models"
)

func TestUpdateSection(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/sections/12345" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPut {
			t.Errorf("unexpected method: %s", r.Method)
		}

		var req struct {
			Data struct {
				Name string `json:"name"`
			} `json:"data"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}

		if req.Data.Name != "Updated Section" {
			t.Errorf("unexpected name: %s", req.Data.Name)
		}

		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"gid":  "12345",
				"name": "Updated Section",
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

	section, err := client.UpdateSection(context.Background(), "12345", models.SectionUpdateRequest{
		Name: "Updated Section",
	})
	if err != nil {
		t.Fatalf("UpdateSection() error = %v", err)
	}

	if section.GID != "12345" {
		t.Errorf("section.GID = %q, want %q", section.GID, "12345")
	}
	if section.Name != "Updated Section" {
		t.Errorf("section.Name = %q, want %q", section.Name, "Updated Section")
	}
}

func TestDeleteSection(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/sections/12345" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodDelete {
			t.Errorf("unexpected method: %s", r.Method)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	cfg := &config.Config{
		AccessToken: "test-token",
		Timeout:     5 * time.Second,
	}
	client := NewHTTPClient(cfg, WithBaseURL(server.URL))

	err := client.DeleteSection(context.Background(), "12345")
	if err != nil {
		t.Fatalf("DeleteSection() error = %v", err)
	}
}

func TestInsertSection(t *testing.T) {
	t.Run("insert before", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/projects/proj123/sections/insert" {
				t.Errorf("unexpected path: %s", r.URL.Path)
			}
			if r.Method != http.MethodPost {
				t.Errorf("unexpected method: %s", r.Method)
			}

			var req struct {
				Data models.SectionInsertRequest `json:"data"`
			}
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				t.Fatalf("failed to decode request: %v", err)
			}

			if req.Data.Section != "sect456" {
				t.Errorf("unexpected section: %s", req.Data.Section)
			}
			if req.Data.BeforeSection == nil || *req.Data.BeforeSection != "sect789" {
				t.Errorf("unexpected before_section: %v", req.Data.BeforeSection)
			}
			if req.Data.AfterSection != nil {
				t.Errorf("unexpected after_section: %v", req.Data.AfterSection)
			}

			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		cfg := &config.Config{
			AccessToken: "test-token",
			Timeout:     5 * time.Second,
		}
		client := NewHTTPClient(cfg, WithBaseURL(server.URL))

		beforeSection := "sect789"
		err := client.InsertSection(context.Background(), "proj123", models.SectionInsertRequest{
			Section:       "sect456",
			BeforeSection: &beforeSection,
		})
		if err != nil {
			t.Fatalf("InsertSection() error = %v", err)
		}
	})

	t.Run("insert after", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/projects/proj123/sections/insert" {
				t.Errorf("unexpected path: %s", r.URL.Path)
			}
			if r.Method != http.MethodPost {
				t.Errorf("unexpected method: %s", r.Method)
			}

			var req struct {
				Data models.SectionInsertRequest `json:"data"`
			}
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				t.Fatalf("failed to decode request: %v", err)
			}

			if req.Data.Section != "sect456" {
				t.Errorf("unexpected section: %s", req.Data.Section)
			}
			if req.Data.BeforeSection != nil {
				t.Errorf("unexpected before_section: %v", req.Data.BeforeSection)
			}
			if req.Data.AfterSection == nil || *req.Data.AfterSection != "sect789" {
				t.Errorf("unexpected after_section: %v", req.Data.AfterSection)
			}

			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		cfg := &config.Config{
			AccessToken: "test-token",
			Timeout:     5 * time.Second,
		}
		client := NewHTTPClient(cfg, WithBaseURL(server.URL))

		afterSection := "sect789"
		err := client.InsertSection(context.Background(), "proj123", models.SectionInsertRequest{
			Section:      "sect456",
			AfterSection: &afterSection,
		})
		if err != nil {
			t.Fatalf("InsertSection() error = %v", err)
		}
	})
}
