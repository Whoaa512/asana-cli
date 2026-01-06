package session

import (
	"os/exec"
	"path/filepath"
	"strings"
)

func GetCurrentBranch() string {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func GetRepoRoot() string {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func GetRepoName() string {
	root := GetRepoRoot()
	if root == "" {
		return ""
	}

	cmd := exec.Command("git", "remote", "get-url", "origin")
	out, err := cmd.Output()
	if err != nil {
		return filepath.Base(root)
	}

	url := strings.TrimSpace(string(out))
	return normalizeRepoURL(url)
}

func normalizeRepoURL(url string) string {
	url = strings.TrimSuffix(url, ".git")

	if strings.HasPrefix(url, "git@") {
		url = strings.TrimPrefix(url, "git@")
		url = strings.Replace(url, ":", "/", 1)
	}

	url = strings.TrimPrefix(url, "https://")
	url = strings.TrimPrefix(url, "http://")

	return url
}

func IsInGitRepo() bool {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	out, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(out)) == "true"
}
