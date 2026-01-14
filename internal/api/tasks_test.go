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

func TestAddTag(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/tasks/12345/addTag" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("unexpected method: %s", r.Method)
		}

		var req struct {
			Data struct {
				Tag string `json:"tag"`
			} `json:"data"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}

		if req.Data.Tag != "tag1" {
			t.Errorf("unexpected tag: %s", req.Data.Tag)
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

	task, err := client.AddTag(context.Background(), "12345", "tag1")
	if err != nil {
		t.Fatalf("AddTag() error = %v", err)
	}

	if task.GID != "12345" {
		t.Errorf("task.GID = %q, want %q", task.GID, "12345")
	}
	if task.Name != "Test Task" {
		t.Errorf("task.Name = %q, want %q", task.Name, "Test Task")
	}
}

func TestRemoveTag(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/tasks/12345/removeTag" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("unexpected method: %s", r.Method)
		}

		var req struct {
			Data struct {
				Tag string `json:"tag"`
			} `json:"data"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}

		if req.Data.Tag != "tag1" {
			t.Errorf("unexpected tag: %s", req.Data.Tag)
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

	task, err := client.RemoveTag(context.Background(), "12345", "tag1")
	if err != nil {
		t.Fatalf("RemoveTag() error = %v", err)
	}

	if task.GID != "12345" {
		t.Errorf("task.GID = %q, want %q", task.GID, "12345")
	}
	if task.Name != "Test Task" {
		t.Errorf("task.Name = %q, want %q", task.Name, "Test Task")
	}
}

func TestAddToProject(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/tasks/12345/addProject" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("unexpected method: %s", r.Method)
		}

		var req struct {
			Data struct {
				Project string `json:"project"`
			} `json:"data"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}

		if req.Data.Project != "proj1" {
			t.Errorf("unexpected project: %s", req.Data.Project)
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

	task, err := client.AddToProject(context.Background(), "12345", "proj1")
	if err != nil {
		t.Fatalf("AddToProject() error = %v", err)
	}

	if task.GID != "12345" {
		t.Errorf("task.GID = %q, want %q", task.GID, "12345")
	}
	if task.Name != "Test Task" {
		t.Errorf("task.Name = %q, want %q", task.Name, "Test Task")
	}
}

func TestRemoveFromProject(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/tasks/12345/removeProject" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("unexpected method: %s", r.Method)
		}

		var req struct {
			Data struct {
				Project string `json:"project"`
			} `json:"data"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}

		if req.Data.Project != "proj1" {
			t.Errorf("unexpected project: %s", req.Data.Project)
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

	task, err := client.RemoveFromProject(context.Background(), "12345", "proj1")
	if err != nil {
		t.Fatalf("RemoveFromProject() error = %v", err)
	}

	if task.GID != "12345" {
		t.Errorf("task.GID = %q, want %q", task.GID, "12345")
	}
	if task.Name != "Test Task" {
		t.Errorf("task.Name = %q, want %q", task.Name, "Test Task")
	}
}

func TestDuplicateTask(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/tasks/12345/duplicate" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("unexpected method: %s", r.Method)
		}

		var req struct {
			Data TaskDuplicateRequest `json:"data"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}

		if req.Data.Name != "Copy of Test Task" {
			t.Errorf("unexpected name: %s", req.Data.Name)
		}
		if len(req.Data.Include) != 2 {
			t.Errorf("expected 2 include options, got %d", len(req.Data.Include))
		}

		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"gid":  "67890",
				"name": "Copy of Test Task",
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

	task, err := client.DuplicateTask(context.Background(), "12345", TaskDuplicateRequest{
		Name:    "Copy of Test Task",
		Include: []string{"subtasks", "attachments"},
	})
	if err != nil {
		t.Fatalf("DuplicateTask() error = %v", err)
	}

	if task.GID != "67890" {
		t.Errorf("task.GID = %q, want %q", task.GID, "67890")
	}
	if task.Name != "Copy of Test Task" {
		t.Errorf("task.Name = %q, want %q", task.Name, "Copy of Test Task")
	}
}

func TestSetParent(t *testing.T) {
	t.Run("set parent", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/tasks/12345/setParent" {
				t.Errorf("unexpected path: %s", r.URL.Path)
			}
			if r.Method != http.MethodPost {
				t.Errorf("unexpected method: %s", r.Method)
			}

			var req struct {
				Data struct {
					Parent *string `json:"parent"`
				} `json:"data"`
			}
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				t.Fatalf("failed to decode request: %v", err)
			}

			if req.Data.Parent == nil || *req.Data.Parent != "99999" {
				t.Errorf("unexpected parent: %v", req.Data.Parent)
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

		parentGID := "99999"
		task, err := client.SetParent(context.Background(), "12345", &parentGID)
		if err != nil {
			t.Fatalf("SetParent() error = %v", err)
		}

		if task.GID != "12345" {
			t.Errorf("task.GID = %q, want %q", task.GID, "12345")
		}
		if task.Name != "Test Task" {
			t.Errorf("task.Name = %q, want %q", task.Name, "Test Task")
		}
	})

	t.Run("clear parent", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/tasks/12345/setParent" {
				t.Errorf("unexpected path: %s", r.URL.Path)
			}
			if r.Method != http.MethodPost {
				t.Errorf("unexpected method: %s", r.Method)
			}

			var req struct {
				Data struct {
					Parent *string `json:"parent"`
				} `json:"data"`
			}
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				t.Fatalf("failed to decode request: %v", err)
			}

			if req.Data.Parent != nil {
				t.Errorf("expected nil parent, got: %v", req.Data.Parent)
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

		task, err := client.SetParent(context.Background(), "12345", nil)
		if err != nil {
			t.Fatalf("SetParent() error = %v", err)
		}

		if task.GID != "12345" {
			t.Errorf("task.GID = %q, want %q", task.GID, "12345")
		}
		if task.Name != "Test Task" {
			t.Errorf("task.Name = %q, want %q", task.Name, "Test Task")
		}
	})
}

func TestListTaskProjects(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/tasks/12345/projects" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("unexpected method: %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{
				{
					"gid":  "proj1",
					"name": "Project One",
				},
				{
					"gid":  "proj2",
					"name": "Project Two",
				},
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

	projects, err := client.ListTaskProjects(context.Background(), "12345")
	if err != nil {
		t.Fatalf("ListTaskProjects() error = %v", err)
	}

	if len(projects) != 2 {
		t.Fatalf("expected 2 projects, got %d", len(projects))
	}

	if projects[0].GID != "proj1" {
		t.Errorf("projects[0].GID = %q, want %q", projects[0].GID, "proj1")
	}
	if projects[0].Name != "Project One" {
		t.Errorf("projects[0].Name = %q, want %q", projects[0].Name, "Project One")
	}

	if projects[1].GID != "proj2" {
		t.Errorf("projects[1].GID = %q, want %q", projects[1].GID, "proj2")
	}
	if projects[1].Name != "Project Two" {
		t.Errorf("projects[1].Name = %q, want %q", projects[1].Name, "Project Two")
	}
}
