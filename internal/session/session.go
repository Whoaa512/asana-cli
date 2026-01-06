package session

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const SessionDir = ".asana-cli"
const SessionFile = "session.json"

type LogEntry struct {
	Timestamp time.Time `json:"ts"`
	Type      string    `json:"type"`
	Text      string    `json:"text"`
}

type Session struct {
	TaskGID     string     `json:"task_gid"`
	ProjectGID  string     `json:"project_gid,omitempty"`
	StartedAt   time.Time  `json:"started_at"`
	Repo        string     `json:"repo,omitempty"`
	StartBranch string     `json:"start_branch,omitempty"`
	Logs        []LogEntry `json:"logs,omitempty"`
	path        string
}

func (s *Session) Path() string {
	return s.path
}

func (s *Session) Duration() time.Duration {
	return time.Since(s.StartedAt)
}

func (s *Session) FormatDuration() string {
	d := s.Duration()
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60

	if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	}
	return fmt.Sprintf("%dm", minutes)
}

func (s *Session) AddLog(logType, text string) {
	s.Logs = append(s.Logs, LogEntry{
		Timestamp: time.Now().UTC(),
		Type:      logType,
		Text:      text,
	})
}

func (s *Session) IsStale() bool {
	return time.Since(s.StartedAt) > 24*time.Hour
}

func Load(dir string) (*Session, error) {
	path := filepath.Join(dir, SessionDir, SessionFile)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var s Session
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, err
	}
	s.path = path

	return &s, nil
}

func (s *Session) Save(dir string) error {
	sessionDir := filepath.Join(dir, SessionDir)
	if err := os.MkdirAll(sessionDir, 0755); err != nil {
		return err
	}
	path := filepath.Join(sessionDir, SessionFile)
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	s.path = path
	return os.WriteFile(path, data, 0644)
}

func Delete(dir string) error {
	path := filepath.Join(dir, SessionDir, SessionFile)
	err := os.Remove(path)
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

func New(taskGID string, opts ...SessionOption) *Session {
	s := &Session{
		TaskGID:   taskGID,
		StartedAt: time.Now().UTC(),
		Logs:      []LogEntry{},
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

type SessionOption func(*Session)

func WithProject(gid string) SessionOption {
	return func(s *Session) {
		s.ProjectGID = gid
	}
}

func WithRepo(repo string) SessionOption {
	return func(s *Session) {
		s.Repo = repo
	}
}

func WithBranch(branch string) SessionOption {
	return func(s *Session) {
		s.StartBranch = branch
	}
}

func (s *Session) FormatSummary(endBranch, extraSummary string) string {
	var result string

	result = "## Work Session\n\n"
	result += fmt.Sprintf("**Duration:** %s\n", s.FormatDuration())

	if s.StartBranch != "" {
		if endBranch != "" && endBranch != s.StartBranch {
			result += fmt.Sprintf("**Branch:** %s → %s\n", s.StartBranch, endBranch)
		} else {
			result += fmt.Sprintf("**Branch:** %s\n", s.StartBranch)
		}
	}

	if s.Repo != "" {
		result += fmt.Sprintf("**Repo:** %s\n", s.Repo)
	}

	progress := s.logsByType("progress")
	decisions := s.logsByType("decision")
	blockers := s.logsByType("blocker")

	if len(progress) > 0 {
		result += "\n### Progress\n"
		for _, l := range progress {
			result += fmt.Sprintf("- %s\n", l.Text)
		}
	}

	if len(decisions) > 0 {
		result += "\n### Decisions\n"
		for _, l := range decisions {
			result += fmt.Sprintf("- %s\n", l.Text)
		}
	}

	if len(blockers) > 0 {
		result += "\n### Blockers\n"
		for _, l := range blockers {
			result += fmt.Sprintf("- %s\n", l.Text)
		}
	}

	if extraSummary != "" {
		result += fmt.Sprintf("\n### Summary\n%s\n", extraSummary)
	}

	result += "\n---\n*Posted via asana-cli*"

	return result
}

func (s *Session) logsByType(logType string) []LogEntry {
	var result []LogEntry
	for _, l := range s.Logs {
		if l.Type == logType {
			result = append(result, l)
		}
	}
	return result
}

func (s *Session) HasLogs() bool {
	return len(s.Logs) > 0
}
