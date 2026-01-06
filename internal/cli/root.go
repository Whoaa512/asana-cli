package cli

import (
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/whoaa512/asana-cli/internal/api"
	"github.com/whoaa512/asana-cli/internal/config"
	"github.com/whoaa512/asana-cli/internal/errors"
	"github.com/whoaa512/asana-cli/internal/output"
)

var (
	flagWorkspace  string
	flagDebug      bool
	flagDryRun     bool
	flagTimeout    time.Duration
	flagConfigPath string
)

var rootCmd = &cobra.Command{
	Use:           "asana",
	Short:         "CLI for interacting with Asana",
	Long:          "A CLI tool for managing Asana tasks, designed for AI agents with JSON-only output.",
	SilenceUsage:  true,
	SilenceErrors: true,
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&flagWorkspace, "workspace", "w", "", "Override workspace GID")
	rootCmd.PersistentFlags().BoolVar(&flagDebug, "debug", false, "Print HTTP requests/responses to stderr")
	rootCmd.PersistentFlags().BoolVar(&flagDryRun, "dry-run", false, "Preview mutations without executing")
	rootCmd.PersistentFlags().DurationVar(&flagTimeout, "timeout", 0, "HTTP request timeout (default 30s)")
	rootCmd.PersistentFlags().StringVar(&flagConfigPath, "config", "", "Config file path (default ~/.config/asana-cli/config.json)")
}

func Execute() int {
	if err := rootCmd.Execute(); err != nil {
		out := output.NewJSON(os.Stdout)
		_ = out.PrintError(err)
		return errors.GetExitCode(err)
	}
	return errors.ExitSuccess
}

func loadConfig() (*config.Config, error) {
	flags := &config.Flags{
		Workspace:  flagWorkspace,
		Debug:      flagDebug,
		DryRun:     flagDryRun,
		Timeout:    flagTimeout,
		ConfigPath: flagConfigPath,
	}
	return config.Load(flags)
}

func newClient(cfg *config.Config) api.Client {
	var opts []api.Option
	if cfg.Debug {
		opts = append(opts, api.WithDebug(os.Stderr))
	}
	return api.NewHTTPClient(cfg, opts...)
}

func requireAuth(cfg *config.Config) error {
	if cfg.AccessToken == "" {
		return errors.NewAuthError("ASANA_ACCESS_TOKEN environment variable not set")
	}
	return nil
}
