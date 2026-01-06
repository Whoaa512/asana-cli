package cli

import (
	"context"
	"os"

	"github.com/spf13/cobra"

	"github.com/whoaa512/asana-cli/internal/output"
)

var taskSubtaskCmd = &cobra.Command{
	Use:   "subtask",
	Short: "Manage task subtasks",
}

var taskSubtaskListCmd = &cobra.Command{
	Use:   "list <task_gid>",
	Short: "List subtasks of a task",
	Args:  cobra.ExactArgs(1),
	RunE:  runTaskSubtaskList,
}

var taskSubtaskAddCmd = &cobra.Command{
	Use:   "add <task_gid>",
	Short: "Add a subtask to a task",
	Args:  cobra.ExactArgs(1),
	RunE:  runTaskSubtaskAdd,
}

var (
	subtaskListLimit  int
	subtaskListOffset string
	subtaskAddName    string
)

func init() {
	taskCmd.AddCommand(taskSubtaskCmd)
	taskSubtaskCmd.AddCommand(taskSubtaskListCmd)
	taskSubtaskCmd.AddCommand(taskSubtaskAddCmd)

	taskSubtaskListCmd.Flags().IntVar(&subtaskListLimit, "limit", 50, "Max results to return")
	taskSubtaskListCmd.Flags().StringVar(&subtaskListOffset, "offset", "", "Pagination offset")

	taskSubtaskAddCmd.Flags().StringVar(&subtaskAddName, "name", "", "Subtask name (required)")
	_ = taskSubtaskAddCmd.MarkFlagRequired("name")
}

func runTaskSubtaskList(_ *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	if err := requireAuth(cfg); err != nil {
		return err
	}

	client := newClient(cfg)
	result, err := client.ListSubtasks(context.Background(), args[0], subtaskListLimit, subtaskListOffset)
	if err != nil {
		return err
	}

	out := output.NewJSON(os.Stdout)
	return out.Print(result)
}

func runTaskSubtaskAdd(_ *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	if err := requireAuth(cfg); err != nil {
		return err
	}

	if cfg.DryRun {
		out := output.NewJSON(os.Stdout)
		return out.Print(map[string]any{"dry_run": true, "parent_gid": args[0], "name": subtaskAddName})
	}

	client := newClient(cfg)
	task, err := client.AddSubtask(context.Background(), args[0], subtaskAddName)
	if err != nil {
		return err
	}

	out := output.NewJSON(os.Stdout)
	return out.Print(task)
}
