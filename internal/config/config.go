package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const (
	DefaultTimeout    = 30 * time.Second
	DefaultConfigPath = "~/.config/asana-cli/config.json"
)

type Config struct {
	AccessToken      string        `json:"-"`
	Workspace        string        `json:"default_workspace,omitempty"`
	Team             string        `json:"default_team,omitempty"`
	Project          string        `json:"-"`
	Task             string        `json:"-"`
	Timeout          time.Duration `json:"-"`
	TimeoutStr       string        `json:"timeout,omitempty"`
	Debug            bool          `json:"debug,omitempty"`
	DryRun           bool          `json:"-"`
	ConfigPath       string        `json:"-"`
	LocalContextPath string        `json:"-"`
	configFileLoaded bool
}

type Flags struct {
	Workspace  string
	Debug      bool
	DryRun     bool
	Timeout    time.Duration
	ConfigPath string
}

func Load(flags *Flags) (*Config, error) {
	cfg := &Config{
		Timeout: DefaultTimeout,
	}

	configPath := expandPath(DefaultConfigPath)
	if flags != nil && flags.ConfigPath != "" {
		configPath = expandPath(flags.ConfigPath)
	}
	cfg.ConfigPath = configPath

	if err := cfg.loadFromFile(configPath); err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	if err := cfg.loadFromLocalContext(); err != nil {
		return nil, err
	}

	cfg.loadFromEnv()

	if flags != nil {
		cfg.applyFlags(flags)
	}

	return cfg, nil
}

func (c *Config) loadFromLocalContext() error {
	ctx, err := LoadLocalContext()
	if err != nil {
		return err
	}

	c.LocalContextPath = ctx.Path()

	if ctx.Workspace != "" {
		c.Workspace = ctx.Workspace
	}
	if ctx.Project != "" {
		c.Project = ctx.Project
	}
	if ctx.Task != "" {
		c.Task = ctx.Task
	}
	return nil
}

func (c *Config) loadFromFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var fileConfig struct {
		DefaultWorkspace string `json:"default_workspace"`
		DefaultTeam      string `json:"default_team"`
		Timeout          string `json:"timeout"`
		Debug            bool   `json:"debug"`
	}

	if err := json.Unmarshal(data, &fileConfig); err != nil {
		return err
	}

	if fileConfig.DefaultWorkspace != "" {
		c.Workspace = fileConfig.DefaultWorkspace
	}
	if fileConfig.DefaultTeam != "" {
		c.Team = fileConfig.DefaultTeam
	}
	if fileConfig.Timeout != "" {
		if d, err := time.ParseDuration(fileConfig.Timeout); err == nil {
			c.Timeout = d
		}
	}
	if fileConfig.Debug {
		c.Debug = true
	}
	c.configFileLoaded = true
	return nil
}

func (c *Config) loadFromEnv() {
	if token := os.Getenv("ASANA_ACCESS_TOKEN"); token != "" {
		c.AccessToken = token
	}
	if ws := os.Getenv("ASANA_WORKSPACE"); ws != "" {
		c.Workspace = ws
	}
	if debug := os.Getenv("ASANA_DEBUG"); debug != "" {
		c.Debug = parseBool(debug)
	}
}

func (c *Config) applyFlags(flags *Flags) {
	if flags.Workspace != "" {
		c.Workspace = flags.Workspace
	}
	if flags.Debug {
		c.Debug = true
	}
	if flags.DryRun {
		c.DryRun = true
	}
	if flags.Timeout > 0 {
		c.Timeout = flags.Timeout
	}
}

func expandPath(path string) string {
	return ExpandPath(path)
}

func ExpandPath(path string) string {
	if len(path) > 0 && path[0] == '~' {
		home, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		return filepath.Join(home, path[1:])
	}
	return path
}

func (c *Config) ConfigFileLoaded() bool {
	return c.configFileLoaded
}

func parseBool(s string) bool {
	if s == "1" || s == "true" || s == "TRUE" || s == "True" {
		return true
	}
	b, _ := strconv.ParseBool(s)
	return b
}
