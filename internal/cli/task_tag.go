package cli

import (
	"context"
	"os"

	"github.com/spf13/cobra"

	"github.com/whoaa512/asana-cli/internal/output"
)

var taskTagCmd = &cobra.Command{
	Use:   "tag",
	Short: "Manage task tags",
}

var taskTagAddCmd = &cobra.Command{
	Use:   "add <task_gid> <tag_gid>",
	Short: "Add a tag to a task",
	Args:  cobra.ExactArgs(2),
	RunE:  runTaskTagAdd,
}

var taskTagRmCmd = &cobra.Command{
	Use:   "rm <task_gid> <tag_gid>",
	Short: "Remove a tag from a task",
	Args:  cobra.ExactArgs(2),
	RunE:  runTaskTagRm,
}

func init() {
	taskCmd.AddCommand(taskTagCmd)
	taskTagCmd.AddCommand(taskTagAddCmd)
	taskTagCmd.AddCommand(taskTagRmCmd)
}

func runTaskTagAdd(_ *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	if err := requireAuth(cfg); err != nil {
		return err
	}

	taskGID := args[0]
	tagGID := args[1]

	if cfg.DryRun {
		out := output.NewJSON(os.Stdout)
		return out.Print(map[string]any{
			"dry_run":  true,
			"task_gid": taskGID,
			"tag_gid":  tagGID,
		})
	}

	client := newClient(cfg)
	task, err := client.AddTag(context.Background(), taskGID, tagGID)
	if err != nil {
		return err
	}

	out := output.NewJSON(os.Stdout)
	return out.Print(task)
}

func runTaskTagRm(_ *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	if err := requireAuth(cfg); err != nil {
		return err
	}

	taskGID := args[0]
	tagGID := args[1]

	if cfg.DryRun {
		out := output.NewJSON(os.Stdout)
		return out.Print(map[string]any{
			"dry_run":  true,
			"task_gid": taskGID,
			"tag_gid":  tagGID,
			"action":   "remove",
		})
	}

	client := newClient(cfg)
	task, err := client.RemoveTag(context.Background(), taskGID, tagGID)
	if err != nil {
		return err
	}

	out := output.NewJSON(os.Stdout)
	return out.Print(task)
}
