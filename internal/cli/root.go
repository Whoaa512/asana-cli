package cli

import (
	"context"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/whoaa512/asana-cli/internal/api"
	"github.com/whoaa512/asana-cli/internal/config"
	"github.com/whoaa512/asana-cli/internal/errors"
	"github.com/whoaa512/asana-cli/internal/models"
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
	Use:   "asana",
	Short: "CLI for interacting with Asana",
	Long: `A CLI tool for managing Asana tasks, designed for AI agents with JSON-only output.

Set ASANA_ACCESS_TOKEN environment variable to authenticate.
Use .asana.json in your repo root to set project/task context.

Global flags:
  --debug     Print HTTP requests/responses to stderr
  --dry-run   Preview mutations without executing
  --workspace Override workspace GID`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&flagWorkspace, "workspace", "w", "", "Override workspace GID")
	rootCmd.PersistentFlags().BoolVar(&flagDebug, "debug", false, "Print HTTP requests/responses to stderr")
	rootCmd.PersistentFlags().BoolVar(&flagDryRun, "dry-run", false, "Preview mutations without executing")
	rootCmd.PersistentFlags().DurationVar(&flagTimeout, "timeout", 0, "HTTP request timeout (default 30s)")
	rootCmd.PersistentFlags().StringVar(&flagConfigPath, "config", "", "Config file path (default ~/.config/asana-cli/config.json)")

	rootCmd.AddCommand(logCmd)
	rootCmd.AddCommand(noteCmd)
	rootCmd.AddCommand(doneCmd)
}

var logCmd = &cobra.Command{
	Use:   "log <message>",
	Short: "Add note to current session (alias for 'session log')",
	Args:  cobra.ExactArgs(1),
	RunE:  runSessionLog,
}

var noteCmd = &cobra.Command{
	Use:   "note <message>",
	Short: "Add comment to context task (alias for 'task comment add')",
	Args:  cobra.ExactArgs(1),
	RunE:  runNote,
}

var doneCmd = &cobra.Command{
	Use:   "done",
	Short: "Mark context task as complete (alias for 'task complete')",
	RunE:  runDone,
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

func runNote(_ *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	if err := requireAuth(cfg); err != nil {
		return err
	}

	if cfg.Task == "" {
		return errors.NewGeneralError("no task in context, set via 'ctx task <gid>'", nil)
	}

	if cfg.DryRun {
		out := output.NewJSON(os.Stdout)
		return out.Print(map[string]any{"dry_run": true, "task_gid": cfg.Task, "text": args[0]})
	}

	client := newClient(cfg)
	story, err := client.AddComment(context.Background(), cfg.Task, args[0])
	if err != nil {
		return err
	}

	out := output.NewJSON(os.Stdout)
	return out.Print(story)
}

func runDone(_ *cobra.Command, _ []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	if err := requireAuth(cfg); err != nil {
		return err
	}

	if cfg.Task == "" {
		return errors.NewGeneralError("no task in context, set via 'ctx task <gid>'", nil)
	}

	completed := true
	req := models.TaskUpdateRequest{Completed: &completed}

	if cfg.DryRun {
		out := output.NewJSON(os.Stdout)
		return out.Print(map[string]any{"dry_run": true, "task_gid": cfg.Task, "action": "complete"})
	}

	client := newClient(cfg)
	task, err := client.UpdateTask(context.Background(), cfg.Task, req)
	if err != nil {
		return err
	}

	out := output.NewJSON(os.Stdout)
	return out.Print(task)
}
