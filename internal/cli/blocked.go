package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

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

	incompleteTasks, err := fetchIncompleteTasksWithDeps(client, project, blockedAssignee, blockedLimit)
	if err != nil {
		return err
	}

	blockedTasks, err := filterBlockedTasks(incompleteTasks)
	if err != nil {
		return err
	}

	out := newOutput()
	return out.PrintTasks(blockedTasks)
}

func filterBlockedTasks(tasks []models.Task) ([]models.Task, error) {
	var blocked []models.Task
	for _, task := range tasks {
		if task.Dependencies == nil {
			return nil, fmt.Errorf("task %s missing dependency data - ensure opt_fields includes dependencies", task.GID)
		}
		if !isTaskReady(task) {
			blocked = append(blocked, task)
		}
	}
	return blocked, nil
}
