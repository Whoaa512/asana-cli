package cli

import (
	"context"
	"os"

	"github.com/spf13/cobra"

	"github.com/whoaa512/asana-cli/internal/api"
	"github.com/whoaa512/asana-cli/internal/errors"
	"github.com/whoaa512/asana-cli/internal/models"
	"github.com/whoaa512/asana-cli/internal/output"
)

var blockedCmd = &cobra.Command{
	Use:   "blocked",
	Short: "List tasks blocked by incomplete dependencies",
	Long:  "List tasks that are blocked by at least one incomplete dependency. Filters by project and assignee.",
	RunE:  runBlocked,
}

var (
	blockedProject  string
	blockedAssignee string
	blockedLimit    int
)

func init() {
	rootCmd.AddCommand(blockedCmd)
	blockedCmd.Flags().StringVar(&blockedProject, "project", "", "Filter by project GID")
	blockedCmd.Flags().StringVar(&blockedAssignee, "assignee", "", "Filter by assignee GID or 'me'")
	blockedCmd.Flags().IntVar(&blockedLimit, "limit", 20, "Max results to return")
}

func runBlocked(_ *cobra.Command, _ []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	if err := requireAuth(cfg); err != nil {
		return err
	}

	project := blockedProject
	if project == "" && cfg.Project != "" {
		project = cfg.Project
	}
	if project == "" {
		return errors.NewGeneralError("no project specified via --project or context", nil)
	}

	if cfg.DryRun {
		out := output.NewJSON(os.Stdout)
		return out.Print(map[string]any{
			"dry_run":  true,
			"project":  project,
			"assignee": blockedAssignee,
			"limit":    blockedLimit,
			"action":   "blocked",
		})
	}

	client := newClient(cfg)

	incompleteTasks, err := fetchIncompleteTasks(client, project, blockedAssignee, blockedLimit)
	if err != nil {
		return err
	}

	blockedTasks, err := filterBlockedTasks(client, incompleteTasks)
	if err != nil {
		return err
	}

	out := output.NewJSON(os.Stdout)
	return out.Print(map[string]any{"data": blockedTasks})
}

func filterBlockedTasks(client api.Client, tasks []models.Task) ([]models.Task, error) {
	var blocked []models.Task

	for _, task := range tasks {
		hasBlockers, err := hasIncompleteDeps(client, task.GID)
		if err != nil {
			return nil, err
		}

		if hasBlockers {
			blocked = append(blocked, task)
		}
	}

	return blocked, nil
}

func hasIncompleteDeps(client api.Client, taskGID string) (bool, error) {
	deps, err := client.ListDependencies(context.Background(), taskGID)
	if err != nil {
		return false, err
	}

	for _, dep := range deps {
		if !dep.Completed {
			return true, nil
		}
	}

	return false, nil
}
