package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const LocalContextFile = ".asana.json"

type LocalContext struct {
	Workspace string `json:"workspace,omitempty"`
	Project   string `json:"project,omitempty"`
	Task      string `json:"task,omitempty"`
	path      string
}

func (lc *LocalContext) Path() string {
	return lc.path
}

func LoadLocalContext() (*LocalContext, error) {
	path, err := findContextFile()
	if err != nil {
		return nil, err
	}
	if path == "" {
		return &LocalContext{}, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var ctx LocalContext
	if err := json.Unmarshal(data, &ctx); err != nil {
		return nil, err
	}
	ctx.path = path

	return &ctx, nil
}

func findContextFile() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	home, err := os.UserHomeDir()
	if err != nil {
		home = ""
	}

	dir := cwd
	for {
		contextPath := filepath.Join(dir, LocalContextFile)
		if _, err := os.Stat(contextPath); err == nil {
			return contextPath, nil
		}

		if isGitRoot(dir) {
			return "", nil
		}

		if dir == home {
			return "", nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", nil
		}
		dir = parent
	}
}

func isGitRoot(dir string) bool {
	gitPath := filepath.Join(dir, ".git")
	info, err := os.Stat(gitPath)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func FindContextFileDir() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	home, err := os.UserHomeDir()
	if err != nil {
		home = ""
	}

	dir := cwd
	for {
		contextPath := filepath.Join(dir, LocalContextFile)
		if _, err := os.Stat(contextPath); err == nil {
			return dir, nil
		}

		if isGitRoot(dir) {
			return dir, nil
		}

		if dir == home {
			return cwd, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return cwd, nil
		}
		dir = parent
	}
}

func (lc *LocalContext) Save(dir string) error {
	path := filepath.Join(dir, LocalContextFile)
	data, err := json.MarshalIndent(lc, "", "  ")
	if err != nil {
		return err
	}
	lc.path = path
	return os.WriteFile(path, data, 0644)
}
