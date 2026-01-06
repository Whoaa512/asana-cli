package session

import "testing"

func TestNormalizeRepoURL(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"git@github.com:org/repo.git", "github.com/org/repo"},
		{"https://github.com/org/repo.git", "github.com/org/repo"},
		{"https://github.com/org/repo", "github.com/org/repo"},
		{"http://github.com/org/repo.git", "github.com/org/repo"},
		{"git@gitlab.com:user/project.git", "gitlab.com/user/project"},
	}

	for _, tc := range tests {
		result := normalizeRepoURL(tc.input)
		if result != tc.expected {
			t.Errorf("normalizeRepoURL(%q) = %q, expected %q", tc.input, result, tc.expected)
		}
	}
}

func TestIsInGitRepo(t *testing.T) {
	result := IsInGitRepo()
	if !result {
		t.Skip("test requires being in a git repo")
	}
}

func TestGetCurrentBranch(t *testing.T) {
	if !IsInGitRepo() {
		t.Skip("test requires being in a git repo")
	}

	branch := GetCurrentBranch()
	if branch == "" {
		t.Error("expected non-empty branch name")
	}
}

func TestGetRepoRoot(t *testing.T) {
	if !IsInGitRepo() {
		t.Skip("test requires being in a git repo")
	}

	root := GetRepoRoot()
	if root == "" {
		t.Error("expected non-empty repo root")
	}
}
