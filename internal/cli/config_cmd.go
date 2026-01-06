package cli

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/whoaa512/asana-cli/internal/config"
	"github.com/whoaa512/asana-cli/internal/errors"
	"github.com/whoaa512/asana-cli/internal/output"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage CLI configuration",
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Display current resolved configuration",
	RunE:  runConfigShow,
}

var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Create global config file",
	Long: `Create a global config file at ~/.config/asana-cli/config.json.

Use flags for non-interactive setup, or run without flags for interactive mode.`,
	RunE: runConfigInit,
}

var (
	configInitWorkspace string
	configInitTimeout   string
	configInitDebug     bool
	configInitForce     bool
)

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configInitCmd)

	configInitCmd.Flags().StringVar(&configInitWorkspace, "workspace", "", "Default workspace GID")
	configInitCmd.Flags().StringVar(&configInitTimeout, "timeout", "", "Request timeout (e.g., 30s, 1m)")
	configInitCmd.Flags().BoolVar(&configInitDebug, "debug", false, "Enable debug mode by default")
	configInitCmd.Flags().BoolVar(&configInitForce, "force", false, "Overwrite existing config file")
}

func runConfigShow(_ *cobra.Command, _ []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	maskedToken := ""
	if cfg.AccessToken != "" {
		if len(cfg.AccessToken) > 8 {
			maskedToken = cfg.AccessToken[:4] + "..." + cfg.AccessToken[len(cfg.AccessToken)-4:]
		} else {
			maskedToken = "****"
		}
	}

	result := map[string]any{
		"access_token":       maskedToken,
		"workspace":          cfg.Workspace,
		"project":            cfg.Project,
		"task":               cfg.Task,
		"timeout":            cfg.Timeout.String(),
		"debug":              cfg.Debug,
		"config_path":        cfg.ConfigPath,
		"config_file_found":  cfg.ConfigFileLoaded(),
		"local_context_path": cfg.LocalContextPath,
	}

	out := output.NewJSON(os.Stdout)
	return out.Print(result)
}

func runConfigInit(_ *cobra.Command, _ []string) error {
	configPath := config.ExpandPath(config.DefaultConfigPath)

	if _, err := os.Stat(configPath); err == nil && !configInitForce {
		return errors.NewGeneralError("config file already exists, use --force to overwrite", nil)
	}

	fileConfig := struct {
		DefaultWorkspace string `json:"default_workspace,omitempty"`
		Timeout          string `json:"timeout,omitempty"`
		Debug            bool   `json:"debug,omitempty"`
	}{}

	hasFlags := configInitWorkspace != "" || configInitTimeout != "" || configInitDebug
	if hasFlags {
		fileConfig.DefaultWorkspace = configInitWorkspace
		fileConfig.Timeout = configInitTimeout
		fileConfig.Debug = configInitDebug
	} else {
		reader := bufio.NewReader(os.Stdin)

		fmt.Fprint(os.Stderr, "Default workspace GID (leave empty to skip): ")
		ws, _ := reader.ReadString('\n')
		fileConfig.DefaultWorkspace = strings.TrimSpace(ws)

		fmt.Fprint(os.Stderr, "Request timeout (e.g., 30s, 1m) [30s]: ")
		timeout, _ := reader.ReadString('\n')
		timeout = strings.TrimSpace(timeout)
		if timeout != "" {
			fileConfig.Timeout = timeout
		}

		fmt.Fprint(os.Stderr, "Enable debug mode by default? (y/N): ")
		debug, _ := reader.ReadString('\n')
		debug = strings.TrimSpace(strings.ToLower(debug))
		fileConfig.Debug = debug == "y" || debug == "yes"
	}

	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return errors.NewGeneralError("failed to create config directory", err)
	}

	data, err := json.MarshalIndent(fileConfig, "", "  ")
	if err != nil {
		return errors.NewGeneralError("failed to encode config", err)
	}

	if err := os.WriteFile(configPath, data, 0600); err != nil {
		return errors.NewGeneralError("failed to write config file", err)
	}

	out := output.NewJSON(os.Stdout)
	return out.Print(map[string]any{
		"created":     true,
		"config_path": configPath,
		"config":      fileConfig,
	})
}
