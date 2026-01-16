package cli

import (
	"testing"

	"github.com/whoaa512/asana-cli/internal/models"
)

func TestFuzzyMatchTasks(t *testing.T) {
	tasks := []models.Task{
		{GID: "1", Name: "Fix login bug"},
		{GID: "2", Name: "Add authentication feature"},
		{GID: "3", Name: "Update user interface"},
		{GID: "4", Name: "Refactor auth module"},
		{GID: "5", Name: "Fix registration bug"},
	}

	tests := []struct {
		name          string
		query         string
		wantGIDs      []string
		wantMinLength int
	}{
		{
			name:          "exact match",
			query:         "Fix login bug",
			wantGIDs:      []string{"1"},
			wantMinLength: 1,
		},
		{
			name:          "fuzzy match auth",
			query:         "auth",
			wantGIDs:      []string{"2", "4"},
			wantMinLength: 2,
		},
		{
			name:          "fuzzy match bug",
			query:         "bug",
			wantGIDs:      []string{"1", "5"},
			wantMinLength: 2,
		},
		{
			name:          "substring match",
			query:         "fix",
			wantGIDs:      []string{"1", "5"},
			wantMinLength: 2,
		},
		{
			name:          "no match",
			query:         "nonexistent",
			wantGIDs:      []string{},
			wantMinLength: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matches := fuzzyMatchTasks(tasks, tt.query)

			if len(matches) < tt.wantMinLength {
				t.Errorf("fuzzyMatchTasks() returned %d matches, want at least %d", len(matches), tt.wantMinLength)
			}

			if len(tt.wantGIDs) > 0 {
				matchedGIDs := make(map[string]bool)
				for _, match := range matches {
					matchedGIDs[match.task.GID] = true
				}

				for _, wantGID := range tt.wantGIDs {
					if !matchedGIDs[wantGID] {
						t.Errorf("fuzzyMatchTasks() missing expected GID %s", wantGID)
					}
				}
			}
		})
	}
}

func TestGIDRegex(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"1234567890", true},
		{"123", true},
		{"0", true},
		{"abc123", false},
		{"123abc", false},
		{"", false},
		{"12 34", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := gidRegex.MatchString(tt.input)
			if got != tt.want {
				t.Errorf("gidRegex.MatchString(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
