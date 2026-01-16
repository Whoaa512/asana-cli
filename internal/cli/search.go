package cli

import (
	"context"
	"os"

	"github.com/spf13/cobra"

	"github.com/whoaa512/asana-cli/internal/api"
	"github.com/whoaa512/asana-cli/internal/errors"
	"github.com/whoaa512/asana-cli/internal/output"
)

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search tasks by text",
	Long: `Search tasks across workspace by text query.

Requires workspace from config. Use --project or --assignee to narrow results.`,
	Example: `  # Search for tasks containing "bug"
  asana search "bug"

  # Search within specific project
  asana search "review" --project 1234567890

  # Search assigned to me
  asana search "urgent" --assignee me

  # Include completed tasks
  asana search "migration" --completed`,
	Args: cobra.ExactArgs(1),
	RunE: runSearch,
}

var (
	searchProject   string
	searchAssignee  string
	searchCompleted bool
	searchLimit     int
	searchOffset    string
)

func init() {
	rootCmd.AddCommand(searchCmd)

	searchCmd.Flags().StringVar(&searchProject, "project", "", "Filter by project GID")
	searchCmd.Flags().StringVar(&searchAssignee, "assignee", "", "Filter by assignee GID or 'me'")
	searchCmd.Flags().BoolVar(&searchCompleted, "completed", false, "Include completed tasks")
	searchCmd.Flags().IntVar(&searchLimit, "limit", 20, "Max results to return")
	searchCmd.Flags().StringVar(&searchOffset, "offset", "", "Pagination offset")
}

func runSearch(_ *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	if err := requireAuth(cfg); err != nil {
		return err
	}

	if cfg.Workspace == "" {
		return errors.NewGeneralError("workspace is required for search", nil)
	}

	opts := api.SearchTasksOptions{
		Workspace: cfg.Workspace,
		Text:      args[0],
		Project:   searchProject,
		Assignee:  searchAssignee,
		Limit:     searchLimit,
		Offset:    searchOffset,
	}

	if !searchCompleted {
		completed := false
		opts.Completed = &completed
	}

	if cfg.DryRun {
		out := output.NewJSON(os.Stdout)
		return out.Print(map[string]any{"dry_run": true, "search_options": opts})
	}

	client := newClient(cfg)
	result, err := client.SearchTasks(context.Background(), opts)
	if err != nil {
		return err
	}

	out := newOutput()
	return out.PrintTaskList(result)
}
