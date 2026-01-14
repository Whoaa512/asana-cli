package cli

import (
	"context"
	"os"

	"github.com/spf13/cobra"

	"github.com/whoaa512/asana-cli/internal/output"
)

var taskProjectCmd = &cobra.Command{
	Use:   "project",
	Short: "Manage task projects",
}

var taskProjectAddCmd = &cobra.Command{
	Use:   "add <task_gid> <project_gid>",
	Short: "Add a task to a project",
	Args:  cobra.ExactArgs(2),
	RunE:  runTaskProjectAdd,
}

var taskProjectRmCmd = &cobra.Command{
	Use:   "rm <task_gid> <project_gid>",
	Short: "Remove a task from a project",
	Args:  cobra.ExactArgs(2),
	RunE:  runTaskProjectRm,
}

var taskProjectListCmd = &cobra.Command{
	Use:   "list <task_gid>",
	Short: "List projects for a task",
	Args:  cobra.ExactArgs(1),
	RunE:  runTaskProjectList,
}

func init() {
	taskCmd.AddCommand(taskProjectCmd)
	taskProjectCmd.AddCommand(taskProjectAddCmd)
	taskProjectCmd.AddCommand(taskProjectRmCmd)
	taskProjectCmd.AddCommand(taskProjectListCmd)
}

func runTaskProjectAdd(_ *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	if err := requireAuth(cfg); err != nil {
		return err
	}

	taskGID := args[0]
	projectGID := args[1]

	if cfg.DryRun {
		out := output.NewJSON(os.Stdout)
		return out.Print(map[string]any{
			"dry_run":     true,
			"task_gid":    taskGID,
			"project_gid": projectGID,
			"action":      "add",
		})
	}

	client := newClient(cfg)
	task, err := client.AddToProject(context.Background(), taskGID, projectGID)
	if err != nil {
		return err
	}

	out := output.NewJSON(os.Stdout)
	return out.Print(task)
}

func runTaskProjectRm(_ *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	if err := requireAuth(cfg); err != nil {
		return err
	}

	taskGID := args[0]
	projectGID := args[1]

	if cfg.DryRun {
		out := output.NewJSON(os.Stdout)
		return out.Print(map[string]any{
			"dry_run":     true,
			"task_gid":    taskGID,
			"project_gid": projectGID,
			"action":      "remove",
		})
	}

	client := newClient(cfg)
	task, err := client.RemoveFromProject(context.Background(), taskGID, projectGID)
	if err != nil {
		return err
	}

	out := output.NewJSON(os.Stdout)
	return out.Print(task)
}

func runTaskProjectList(_ *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	if err := requireAuth(cfg); err != nil {
		return err
	}

	taskGID := args[0]

	client := newClient(cfg)
	projects, err := client.ListTaskProjects(context.Background(), taskGID)
	if err != nil {
		return err
	}

	out := output.NewJSON(os.Stdout)
	return out.Print(projects)
}
