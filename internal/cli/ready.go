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

var readyCmd = &cobra.Command{
	Use:   "ready",
	Short: "List tasks with no incomplete dependencies",
	Long:  "List tasks that are ready to work on (no blocking dependencies). Filters by project and assignee.",
	RunE:  runReady,
}

var (
	readyProject  string
	readyAssignee string
	readyLimit    int
)

func init() {
	rootCmd.AddCommand(readyCmd)
	readyCmd.Flags().StringVar(&readyProject, "project", "", "Filter by project GID")
	readyCmd.Flags().StringVar(&readyAssignee, "assignee", "", "Filter by assignee GID or 'me'")
	readyCmd.Flags().IntVar(&readyLimit, "limit", 20, "Max results to return")
}

func runReady(_ *cobra.Command, _ []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	if err := requireAuth(cfg); err != nil {
		return err
	}

	project := readyProject
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
			"assignee": readyAssignee,
			"limit":    readyLimit,
			"action":   "ready",
		})
	}

	client := newClient(cfg)

	incompleteTasks, err := fetchIncompleteTasks(client, project, readyAssignee, readyLimit)
	if err != nil {
		return err
	}

	readyTasks, err := filterReadyTasks(client, incompleteTasks)
	if err != nil {
		return err
	}

	out := output.NewJSON(os.Stdout)
	return out.Print(map[string]any{"data": readyTasks})
}

func fetchIncompleteTasks(client api.Client, project, assignee string, limit int) ([]models.Task, error) {
	completed := false
	opts := api.TaskListOptions{
		Project:   project,
		Assignee:  assignee,
		Completed: &completed,
		Limit:     limit,
	}

	result, err := client.ListTasks(context.Background(), opts)
	if err != nil {
		return nil, err
	}

	return result.Data, nil
}

func filterReadyTasks(client api.Client, tasks []models.Task) ([]models.Task, error) {
	var ready []models.Task

	for _, task := range tasks {
		hasBlockers, err := hasIncompleteDependencies(client, task.GID)
		if err != nil {
			return nil, err
		}

		if !hasBlockers {
			ready = append(ready, task)
		}
	}

	return ready, nil
}

func hasIncompleteDependencies(client api.Client, taskGID string) (bool, error) {
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
