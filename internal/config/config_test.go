package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoadDefaults(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv("ASANA_ACCESS_TOKEN", "")
	t.Setenv("ASANA_WORKSPACE", "")
	t.Setenv("ASANA_DEBUG", "")

	cfg, err := Load(nil)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Timeout != DefaultTimeout {
		t.Errorf("Timeout = %v, want %v", cfg.Timeout, DefaultTimeout)
	}
	if cfg.Debug {
		t.Error("Debug should be false by default")
	}
}

func TestLoadFromEnv(t *testing.T) {
	t.Setenv("ASANA_ACCESS_TOKEN", "test-token")
	t.Setenv("ASANA_WORKSPACE", "12345")
	t.Setenv("ASANA_DEBUG", "1")

	cfg, err := Load(nil)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.AccessToken != "test-token" {
		t.Errorf("AccessToken = %q, want %q", cfg.AccessToken, "test-token")
	}
	if cfg.Workspace != "12345" {
		t.Errorf("Workspace = %q, want %q", cfg.Workspace, "12345")
	}
	if !cfg.Debug {
		t.Error("Debug should be true")
	}
}

func TestLoadFromFile(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")
	content := `{"default_workspace": "99999", "timeout": "60s", "debug": true}`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	t.Setenv("ASANA_ACCESS_TOKEN", "")
	t.Setenv("ASANA_WORKSPACE", "")
	t.Setenv("ASANA_DEBUG", "")

	flags := &Flags{ConfigPath: configPath}
	cfg, err := Load(flags)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Workspace != "99999" {
		t.Errorf("Workspace = %q, want %q", cfg.Workspace, "99999")
	}
	if cfg.Timeout != 60*time.Second {
		t.Errorf("Timeout = %v, want %v", cfg.Timeout, 60*time.Second)
	}
	if !cfg.Debug {
		t.Error("Debug should be true from file")
	}
}

func TestFlagsOverrideEnv(t *testing.T) {
	t.Setenv("ASANA_WORKSPACE", "env-workspace")

	flags := &Flags{
		Workspace: "flag-workspace",
		Debug:     true,
		Timeout:   45 * time.Second,
	}

	cfg, err := Load(flags)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Workspace != "flag-workspace" {
		t.Errorf("Workspace = %q, want %q", cfg.Workspace, "flag-workspace")
	}
	if cfg.Timeout != 45*time.Second {
		t.Errorf("Timeout = %v, want %v", cfg.Timeout, 45*time.Second)
	}
}

func TestParseBool(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"1", true},
		{"true", true},
		{"TRUE", true},
		{"True", true},
		{"0", false},
		{"false", false},
		{"", false},
		{"invalid", false},
	}

	for _, tt := range tests {
		if got := parseBool(tt.input); got != tt.want {
			t.Errorf("parseBool(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestExpandPath(t *testing.T) {
	home, _ := os.UserHomeDir()

	tests := []struct {
		input string
		want  string
	}{
		{"~/foo", filepath.Join(home, "foo")},
		{"/absolute/path", "/absolute/path"},
		{"relative/path", "relative/path"},
	}

	for _, tt := range tests {
		if got := expandPath(tt.input); got != tt.want {
			t.Errorf("expandPath(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestLoadIntegratesLocalContext(t *testing.T) {
	tmp := t.TempDir()
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(oldWd) }()

	contextContent := `{"workspace": "local-ws", "project": "local-proj", "task": "local-task"}`
	if err := os.WriteFile(filepath.Join(tmp, LocalContextFile), []byte(contextContent), 0644); err != nil {
		t.Fatal(err)
	}

	if err := os.Chdir(tmp); err != nil {
		t.Fatal(err)
	}

	t.Setenv("ASANA_ACCESS_TOKEN", "")
	t.Setenv("ASANA_WORKSPACE", "")

	cfg, err := Load(nil)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Workspace != "local-ws" {
		t.Errorf("Workspace = %q, want %q", cfg.Workspace, "local-ws")
	}
	if cfg.Project != "local-proj" {
		t.Errorf("Project = %q, want %q", cfg.Project, "local-proj")
	}
	if cfg.Task != "local-task" {
		t.Errorf("Task = %q, want %q", cfg.Task, "local-task")
	}
	if cfg.LocalContextPath == "" {
		t.Error("LocalContextPath should be set")
	}
}

func TestLoadPrecedence_EnvOverridesLocalContext(t *testing.T) {
	tmp := t.TempDir()
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(oldWd) }()

	contextContent := `{"workspace": "local-ws"}`
	if err := os.WriteFile(filepath.Join(tmp, LocalContextFile), []byte(contextContent), 0644); err != nil {
		t.Fatal(err)
	}

	if err := os.Chdir(tmp); err != nil {
		t.Fatal(err)
	}

	t.Setenv("ASANA_WORKSPACE", "env-ws")

	cfg, err := Load(nil)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Workspace != "env-ws" {
		t.Errorf("Workspace = %q, want %q (env should override local)", cfg.Workspace, "env-ws")
	}
}

func TestLoadPrecedence_FlagsOverrideAll(t *testing.T) {
	tmp := t.TempDir()
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(oldWd) }()

	contextContent := `{"workspace": "local-ws"}`
	if err := os.WriteFile(filepath.Join(tmp, LocalContextFile), []byte(contextContent), 0644); err != nil {
		t.Fatal(err)
	}

	if err := os.Chdir(tmp); err != nil {
		t.Fatal(err)
	}

	t.Setenv("ASANA_WORKSPACE", "env-ws")

	flags := &Flags{Workspace: "flag-ws"}
	cfg, err := Load(flags)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Workspace != "flag-ws" {
		t.Errorf("Workspace = %q, want %q (flags should override all)", cfg.Workspace, "flag-ws")
	}
}
