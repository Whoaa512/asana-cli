package session

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewSession(t *testing.T) {
	sess := New("12345")
	if sess.TaskGID != "12345" {
		t.Errorf("expected task_gid 12345, got %s", sess.TaskGID)
	}
	if sess.StartedAt.IsZero() {
		t.Error("expected non-zero started_at")
	}
	if sess.Logs == nil {
		t.Error("expected logs to be initialized")
	}
}

func TestNewSessionWithOptions(t *testing.T) {
	sess := New("12345",
		WithProject("proj123"),
		WithRepo("github.com/user/repo"),
		WithBranch("feature/test"),
	)

	if sess.ProjectGID != "proj123" {
		t.Errorf("expected project_gid proj123, got %s", sess.ProjectGID)
	}
	if sess.Repo != "github.com/user/repo" {
		t.Errorf("expected repo github.com/user/repo, got %s", sess.Repo)
	}
	if sess.StartBranch != "feature/test" {
		t.Errorf("expected branch feature/test, got %s", sess.StartBranch)
	}
}

func TestSessionAddLog(t *testing.T) {
	sess := New("12345")
	sess.AddLog("progress", "did something")
	sess.AddLog("decision", "chose option A")

	if len(sess.Logs) != 2 {
		t.Errorf("expected 2 logs, got %d", len(sess.Logs))
	}
	if sess.Logs[0].Type != "progress" {
		t.Errorf("expected first log type progress, got %s", sess.Logs[0].Type)
	}
	if sess.Logs[0].Text != "did something" {
		t.Errorf("expected first log text 'did something', got %s", sess.Logs[0].Text)
	}
	if sess.Logs[1].Type != "decision" {
		t.Errorf("expected second log type decision, got %s", sess.Logs[1].Type)
	}
}

func TestSessionHasLogs(t *testing.T) {
	sess := New("12345")
	if sess.HasLogs() {
		t.Error("expected HasLogs to be false for new session")
	}

	sess.AddLog("progress", "test")
	if !sess.HasLogs() {
		t.Error("expected HasLogs to be true after adding log")
	}
}

func TestSessionIsStale(t *testing.T) {
	sess := New("12345")
	if sess.IsStale() {
		t.Error("new session should not be stale")
	}

	sess.StartedAt = time.Now().Add(-25 * time.Hour)
	if !sess.IsStale() {
		t.Error("session older than 24h should be stale")
	}
}

func TestSessionSaveAndLoad(t *testing.T) {
	dir := t.TempDir()

	sess := New("12345",
		WithProject("proj123"),
		WithBranch("main"),
	)
	sess.AddLog("progress", "test log")

	if err := sess.Save(dir); err != nil {
		t.Fatalf("failed to save session: %v", err)
	}

	loaded, err := Load(dir)
	if err != nil {
		t.Fatalf("failed to load session: %v", err)
	}
	if loaded == nil {
		t.Fatal("expected loaded session to be non-nil")
	}

	if loaded.TaskGID != "12345" {
		t.Errorf("expected task_gid 12345, got %s", loaded.TaskGID)
	}
	if loaded.ProjectGID != "proj123" {
		t.Errorf("expected project_gid proj123, got %s", loaded.ProjectGID)
	}
	if loaded.StartBranch != "main" {
		t.Errorf("expected branch main, got %s", loaded.StartBranch)
	}
	if len(loaded.Logs) != 1 {
		t.Errorf("expected 1 log, got %d", len(loaded.Logs))
	}
	if loaded.Path() != filepath.Join(dir, SessionFile) {
		t.Errorf("unexpected path: %s", loaded.Path())
	}
}

func TestLoadNonexistent(t *testing.T) {
	dir := t.TempDir()

	sess, err := Load(dir)
	if err != nil {
		t.Fatalf("expected no error for nonexistent session, got %v", err)
	}
	if sess != nil {
		t.Error("expected nil session for nonexistent file")
	}
}

func TestDelete(t *testing.T) {
	dir := t.TempDir()

	sess := New("12345")
	if err := sess.Save(dir); err != nil {
		t.Fatalf("failed to save session: %v", err)
	}

	path := filepath.Join(dir, SessionFile)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatal("session file should exist after save")
	}

	if err := Delete(dir); err != nil {
		t.Fatalf("failed to delete session: %v", err)
	}

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Error("session file should not exist after delete")
	}
}

func TestDeleteNonexistent(t *testing.T) {
	dir := t.TempDir()

	if err := Delete(dir); err != nil {
		t.Errorf("delete of nonexistent should not error, got %v", err)
	}
}

func TestFormatDuration(t *testing.T) {
	sess := New("12345")

	sess.StartedAt = time.Now().Add(-135 * time.Minute)
	formatted := sess.FormatDuration()
	if formatted != "2h 15m" {
		t.Errorf("expected '2h 15m', got '%s'", formatted)
	}

	sess.StartedAt = time.Now().Add(-45 * time.Minute)
	formatted = sess.FormatDuration()
	if formatted != "45m" {
		t.Errorf("expected '45m', got '%s'", formatted)
	}
}

func TestFormatSummary(t *testing.T) {
	sess := New("12345", WithBranch("feature/test"), WithRepo("github.com/org/repo"))
	sess.StartedAt = time.Now().Add(-30 * time.Minute)
	sess.AddLog("progress", "Did task A")
	sess.AddLog("progress", "Did task B")
	sess.AddLog("decision", "Chose approach X")
	sess.AddLog("blocker", "Waiting on API")

	summary := sess.FormatSummary("feature/test", "All done!")

	expected := []string{
		"## Work Session",
		"**Duration:** 30m",
		"**Branch:** feature/test",
		"**Repo:** github.com/org/repo",
		"### Progress",
		"- Did task A",
		"- Did task B",
		"### Decisions",
		"- Chose approach X",
		"### Blockers",
		"- Waiting on API",
		"### Summary",
		"All done!",
		"*Posted via asana-cli*",
	}

	for _, exp := range expected {
		if !containsString(summary, exp) {
			t.Errorf("expected summary to contain '%s'", exp)
		}
	}
}

func TestFormatSummaryBranchChange(t *testing.T) {
	sess := New("12345", WithBranch("main"))
	sess.StartedAt = time.Now().Add(-10 * time.Minute)
	sess.AddLog("progress", "Test")

	summary := sess.FormatSummary("feature/new", "")

	if !containsString(summary, "**Branch:** main → feature/new") {
		t.Error("expected summary to show branch change")
	}
}

func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsSubstring(s, substr))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
