package cli

import (
	"strings"
	"testing"

	"github.com/whoaa512/asana-cli/internal/models"
)

func TestPrimeCommandRegistered(t *testing.T) {
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Use == "prime" {
			found = true
			break
		}
	}
	if !found {
		t.Error("prime command should be registered")
	}
}

func TestPrimeCommandFlags(t *testing.T) {
	flags := []string{"project", "limit", "include-completed"}
	for _, name := range flags {
		if primeCmd.Flags().Lookup(name) == nil {
			t.Errorf("flag %q not registered", name)
		}
	}
}

func TestTruncate(t *testing.T) {
	tests := []struct {
		input  string
		maxLen int
		want   string
	}{
		{"hello", 10, "hello"},
		{"hello world", 5, "hello..."},
		{"", 5, ""},
		{"abc", 3, "abc"},
		{"abcd", 3, "abc..."},
	}

	for _, tt := range tests {
		got := truncate(tt.input, tt.maxLen)
		if got != tt.want {
			t.Errorf("truncate(%q, %d) = %q, want %q", tt.input, tt.maxLen, got, tt.want)
		}
	}
}

func TestWriteCompletedTasks(t *testing.T) {
	t.Run("empty tasks", func(t *testing.T) {
		var output strings.Builder
		writeCompletedTasks(&output, nil)
		if output.String() != "" {
			t.Errorf("expected empty output for nil tasks, got %q", output.String())
		}
	})

	t.Run("with tasks", func(t *testing.T) {
		var output strings.Builder
		tasks := []models.Task{
			{GID: "123", Name: "Task 1", Completed: true, CompletedAt: "2024-01-15T10:00:00Z"},
			{GID: "456", Name: "Task 2", Completed: true, CompletedAt: "2024-01-14T09:00:00Z"},
		}
		writeCompletedTasks(&output, tasks)

		result := output.String()
		if !strings.Contains(result, "## Recently Completed") {
			t.Error("output should contain section header")
		}
		if !strings.Contains(result, "[x] Task 1 (123)") {
			t.Error("output should contain Task 1")
		}
		if !strings.Contains(result, "[x] Task 2 (456)") {
			t.Error("output should contain Task 2")
		}
		if !strings.Contains(result, "completed 2024-01-15") {
			t.Error("output should contain completion date")
		}
	})

	t.Run("without completed_at", func(t *testing.T) {
		var output strings.Builder
		tasks := []models.Task{
			{GID: "789", Name: "Task 3", Completed: true},
		}
		writeCompletedTasks(&output, tasks)

		result := output.String()
		if !strings.Contains(result, "[x] Task 3 (789)") {
			t.Error("output should contain Task 3")
		}
		if strings.Contains(result, "completed") {
			t.Error("output should not contain 'completed' when no date")
		}
	})
}

func TestWriteReadyTasks(t *testing.T) {
	t.Run("empty tasks", func(t *testing.T) {
		var output strings.Builder
		err := writeReadyTasks(&output, nil, nil, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if output.String() != "" {
			t.Errorf("expected empty output for nil tasks, got %q", output.String())
		}
	})
}

func TestWriteBlockedTasks(t *testing.T) {
	t.Run("empty tasks", func(t *testing.T) {
		var output strings.Builder
		err := writeBlockedTasks(&output, nil, nil, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if output.String() != "" {
			t.Errorf("expected empty output for nil tasks, got %q", output.String())
		}
	})
}
