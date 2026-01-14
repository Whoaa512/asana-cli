package cli

import (
	"context"
	"os"

	"github.com/spf13/cobra"

	"github.com/whoaa512/asana-cli/internal/errors"
	"github.com/whoaa512/asana-cli/internal/output"
)

var taskSetParentCmd = &cobra.Command{
	Use:   "set-parent <task_gid> [parent_task_gid]",
	Short: "Set or clear task parent",
	Long: `Set a task's parent (make it a subtask) or clear the parent (make it top-level).

Use --clear flag to remove parent and make the task top-level.`,
	Example: `  # Make task a subtask of another task
  asana task set-parent 1234567890 9876543210

  # Remove parent and make task top-level
  asana task set-parent 1234567890 --clear`,
	Args: cobra.RangeArgs(1, 2),
	RunE: runTaskSetParent,
}

var (
	taskSetParentClear bool
)

func init() {
	taskCmd.AddCommand(taskSetParentCmd)
	taskSetParentCmd.Flags().BoolVar(&taskSetParentClear, "clear", false, "Clear parent (make task top-level)")
}

func runTaskSetParent(_ *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	if err := requireAuth(cfg); err != nil {
		return err
	}

	taskGID := args[0]
	var parentGID *string

	if taskSetParentClear {
		parentGID = nil
	} else {
		if len(args) < 2 {
			return errors.NewGeneralError("parent_task_gid is required when --clear is not set", nil)
		}
		parentGID = &args[1]
	}

	if cfg.DryRun {
		out := output.NewJSON(os.Stdout)
		return out.Print(map[string]any{
			"dry_run":    true,
			"task_gid":   taskGID,
			"parent_gid": parentGID,
		})
	}

	client := newClient(cfg)
	task, err := client.SetParent(context.Background(), taskGID, parentGID)
	if err != nil {
		return err
	}

	out := output.NewJSON(os.Stdout)
	return out.Print(task)
}
