package cli

import (
	"context"
	"os"

	"github.com/spf13/cobra"

	"github.com/whoaa512/asana-cli/internal/api"
	"github.com/whoaa512/asana-cli/internal/output"
)

var taskDuplicateCmd = &cobra.Command{
	Use:   "duplicate <gid>",
	Short: "Duplicate a task",
	Long: `Duplicate an existing task in Asana.

Optionally rename the new task and include specific elements like subtasks, attachments, etc.`,
	Example: `  # Duplicate a task
  asana task duplicate 1234567890

  # Duplicate with a new name
  asana task duplicate 1234567890 --name "Copy of original task"

  # Duplicate including subtasks and attachments
  asana task duplicate 1234567890 --include-subtasks --include-attachments`,
	Args: cobra.ExactArgs(1),
	RunE: runTaskDuplicate,
}

var (
	duplicateName               string
	duplicateIncludeSubtasks    bool
	duplicateIncludeAttachments bool
)

func init() {
	taskCmd.AddCommand(taskDuplicateCmd)

	taskDuplicateCmd.Flags().StringVar(&duplicateName, "name", "", "New task name")
	taskDuplicateCmd.Flags().BoolVar(&duplicateIncludeSubtasks, "include-subtasks", false, "Include subtasks in duplicate")
	taskDuplicateCmd.Flags().BoolVar(&duplicateIncludeAttachments, "include-attachments", false, "Include attachments in duplicate")
}

func runTaskDuplicate(_ *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	if err := requireAuth(cfg); err != nil {
		return err
	}

	taskGID := args[0]
	req := api.TaskDuplicateRequest{
		Name: duplicateName,
	}

	var include []string
	if duplicateIncludeSubtasks {
		include = append(include, "subtasks")
	}
	if duplicateIncludeAttachments {
		include = append(include, "attachments")
	}
	if len(include) > 0 {
		req.Include = include
	}

	if cfg.DryRun {
		out := output.NewJSON(os.Stdout)
		return out.Print(map[string]any{"dry_run": true, "gid": taskGID, "request": req})
	}

	client := newClient(cfg)
	task, err := client.DuplicateTask(context.Background(), taskGID, req)
	if err != nil {
		return err
	}

	out := output.NewJSON(os.Stdout)
	return out.Print(task)
}
