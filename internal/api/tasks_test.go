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

func TestAddFollowers(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/tasks/12345/addFollowers" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("unexpected method: %s", r.Method)
		}

		var req struct {
			Data struct {
				Followers []string `json:"followers"`
			} `json:"data"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}

		if len(req.Data.Followers) != 2 {
			t.Errorf("expected 2 followers, got %d", len(req.Data.Followers))
		}
		if req.Data.Followers[0] != "user1" || req.Data.Followers[1] != "user2" {
			t.Errorf("unexpected followers: %v", req.Data.Followers)
		}

		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"gid":  "12345",
				"name": "Test Task",
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

	task, err := client.AddFollowers(context.Background(), "12345", []string{"user1", "user2"})
	if err != nil {
		t.Fatalf("AddFollowers() error = %v", err)
	}

	if task.GID != "12345" {
		t.Errorf("task.GID = %q, want %q", task.GID, "12345")
	}
	if task.Name != "Test Task" {
		t.Errorf("task.Name = %q, want %q", task.Name, "Test Task")
	}
}

func TestRemoveFollower(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/tasks/12345/removeFollower" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("unexpected method: %s", r.Method)
		}

		var req struct {
			Data struct {
				Follower string `json:"follower"`
			} `json:"data"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}

		if req.Data.Follower != "user1" {
			t.Errorf("unexpected follower: %s", req.Data.Follower)
		}

		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"gid":  "12345",
				"name": "Test Task",
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

	task, err := client.RemoveFollower(context.Background(), "12345", "user1")
	if err != nil {
		t.Fatalf("RemoveFollower() error = %v", err)
	}

	if task.GID != "12345" {
		t.Errorf("task.GID = %q, want %q", task.GID, "12345")
	}
	if task.Name != "Test Task" {
		t.Errorf("task.Name = %q, want %q", task.Name, "Test Task")
	}
}
