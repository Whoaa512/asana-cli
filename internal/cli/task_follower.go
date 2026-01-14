package cli

import (
	"context"
	"os"

	"github.com/spf13/cobra"

	"github.com/whoaa512/asana-cli/internal/output"
)

var taskFollowerCmd = &cobra.Command{
	Use:   "follower",
	Short: "Manage task followers",
}

var taskFollowerAddCmd = &cobra.Command{
	Use:   "add <task_gid> <follower_gid>",
	Short: "Add a follower to a task",
	Args:  cobra.ExactArgs(2),
	RunE:  runTaskFollowerAdd,
}

var taskFollowerRmCmd = &cobra.Command{
	Use:   "rm <task_gid> <follower_gid>",
	Short: "Remove a follower from a task",
	Args:  cobra.ExactArgs(2),
	RunE:  runTaskFollowerRm,
}

func init() {
	taskCmd.AddCommand(taskFollowerCmd)
	taskFollowerCmd.AddCommand(taskFollowerAddCmd)
	taskFollowerCmd.AddCommand(taskFollowerRmCmd)
}

func runTaskFollowerAdd(_ *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	if err := requireAuth(cfg); err != nil {
		return err
	}

	taskGID := args[0]
	followerGID := args[1]

	if cfg.DryRun {
		out := output.NewJSON(os.Stdout)
		return out.Print(map[string]any{
			"dry_run":      true,
			"task_gid":     taskGID,
			"follower_gid": followerGID,
		})
	}

	client := newClient(cfg)
	task, err := client.AddFollowers(context.Background(), taskGID, []string{followerGID})
	if err != nil {
		return err
	}

	out := output.NewJSON(os.Stdout)
	return out.Print(task)
}

func runTaskFollowerRm(_ *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	if err := requireAuth(cfg); err != nil {
		return err
	}

	taskGID := args[0]
	followerGID := args[1]

	if cfg.DryRun {
		out := output.NewJSON(os.Stdout)
		return out.Print(map[string]any{
			"dry_run":      true,
			"task_gid":     taskGID,
			"follower_gid": followerGID,
			"action":       "remove",
		})
	}

	client := newClient(cfg)
	task, err := client.RemoveFollower(context.Background(), taskGID, followerGID)
	if err != nil {
		return err
	}

	out := output.NewJSON(os.Stdout)
	return out.Print(task)
}
