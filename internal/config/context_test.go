package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadLocalContext_NotFound(t *testing.T) {
	tmp := t.TempDir()
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	defer func() {
		if err := os.Chdir(oldWd); err != nil {
			t.Errorf("failed to restore working directory: %v", err)
		}
	}()
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("failed to change to temp directory: %v", err)
	}

	ctx, err := LoadLocalContext()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ctx.Workspace != "" || ctx.Project != "" || ctx.Task != "" {
		t.Error("expected empty context when no file found")
	}
}

func TestLoadLocalContext_Found(t *testing.T) {
	tmp := t.TempDir()
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	defer func() {
		if err := os.Chdir(oldWd); err != nil {
			t.Errorf("failed to restore working directory: %v", err)
		}
	}()

	content := `{"workspace": "123", "project": "456", "task": "789"}`
	if err := os.WriteFile(filepath.Join(tmp, LocalContextFile), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("failed to change to temp directory: %v", err)
	}

	ctx, err := LoadLocalContext()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ctx.Workspace != "123" {
		t.Errorf("expected workspace 123, got %s", ctx.Workspace)
	}
	if ctx.Project != "456" {
		t.Errorf("expected project 456, got %s", ctx.Project)
	}
	if ctx.Task != "789" {
		t.Errorf("expected task 789, got %s", ctx.Task)
	}
}

func TestLoadLocalContext_WalkUp(t *testing.T) {
	tmp := t.TempDir()
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	defer func() {
		if err := os.Chdir(oldWd); err != nil {
			t.Errorf("failed to restore working directory: %v", err)
		}
	}()

	content := `{"workspace": "parent-ws"}`
	if err := os.WriteFile(filepath.Join(tmp, LocalContextFile), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	subdir := filepath.Join(tmp, "subdir", "nested")
	if err := os.MkdirAll(subdir, 0755); err != nil {
		t.Fatal(err)
	}

	if err := os.Chdir(subdir); err != nil {
		t.Fatalf("failed to change to subdir: %v", err)
	}

	ctx, err := LoadLocalContext()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ctx.Workspace != "parent-ws" {
		t.Errorf("expected workspace parent-ws, got %s", ctx.Workspace)
	}
}

func TestLoadLocalContext_StopsAtGitRoot(t *testing.T) {
	tmp := t.TempDir()
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	defer func() {
		if err := os.Chdir(oldWd); err != nil {
			t.Errorf("failed to restore working directory: %v", err)
		}
	}()

	content := `{"workspace": "above-git"}`
	if err := os.WriteFile(filepath.Join(tmp, LocalContextFile), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	gitDir := filepath.Join(tmp, "repo")
	if err := os.MkdirAll(filepath.Join(gitDir, ".git"), 0755); err != nil {
		t.Fatal(err)
	}

	subdir := filepath.Join(gitDir, "src")
	if err := os.MkdirAll(subdir, 0755); err != nil {
		t.Fatal(err)
	}

	if err := os.Chdir(subdir); err != nil {
		t.Fatalf("failed to change to subdir: %v", err)
	}

	ctx, err := LoadLocalContext()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ctx.Workspace != "" {
		t.Errorf("expected empty workspace (stopped at git root), got %s", ctx.Workspace)
	}
}

func TestLocalContext_Save(t *testing.T) {
	tmp := t.TempDir()

	ctx := &LocalContext{
		Workspace: "ws-123",
		Project:   "proj-456",
	}

	if err := ctx.Save(tmp); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(tmp, LocalContextFile))
	if err != nil {
		t.Fatalf("failed to read saved file: %v", err)
	}

	expected := `{
  "workspace": "ws-123",
  "project": "proj-456"
}`
	if string(data) != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, string(data))
	}
}

func TestFindContextFileDir_ReturnsGitRoot(t *testing.T) {
	tmp := t.TempDir()
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	defer func() {
		if err := os.Chdir(oldWd); err != nil {
			t.Errorf("failed to restore working directory: %v", err)
		}
	}()

	gitDir := filepath.Join(tmp, "repo")
	if err := os.MkdirAll(filepath.Join(gitDir, ".git"), 0755); err != nil {
		t.Fatal(err)
	}

	subdir := filepath.Join(gitDir, "src", "pkg")
	if err := os.MkdirAll(subdir, 0755); err != nil {
		t.Fatal(err)
	}

	if err := os.Chdir(subdir); err != nil {
		t.Fatalf("failed to change to subdir: %v", err)
	}

	dir, err := FindContextFileDir()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedReal, _ := filepath.EvalSymlinks(gitDir)
	gotReal, _ := filepath.EvalSymlinks(dir)
	if gotReal != expectedReal {
		t.Errorf("expected git root %s, got %s", expectedReal, gotReal)
	}
}
