package cli

import (
	"context"
	"os"

	"github.com/spf13/cobra"

	"github.com/whoaa512/asana-cli/internal/output"
)

var taskDepCmd = &cobra.Command{
	Use:   "dep",
	Short: "Manage task dependencies",
}

var taskDepAddCmd = &cobra.Command{
	Use:   "add <task_gid> <depends_on_gid>",
	Short: "Add a dependency (mark task as blocked by depends_on)",
	Args:  cobra.ExactArgs(2),
	RunE:  runTaskDepAdd,
}

var taskDepListCmd = &cobra.Command{
	Use:   "list <task_gid>",
	Short: "List task dependencies (both directions)",
	Args:  cobra.ExactArgs(1),
	RunE:  runTaskDepList,
}

var taskDepRmCmd = &cobra.Command{
	Use:   "rm <task_gid> <depends_on_gid>",
	Short: "Remove a dependency",
	Args:  cobra.ExactArgs(2),
	RunE:  runTaskDepRm,
}

func init() {
	taskCmd.AddCommand(taskDepCmd)
	taskDepCmd.AddCommand(taskDepAddCmd)
	taskDepCmd.AddCommand(taskDepListCmd)
	taskDepCmd.AddCommand(taskDepRmCmd)
}

func runTaskDepAdd(_ *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	if err := requireAuth(cfg); err != nil {
		return err
	}

	taskGID := args[0]
	dependsOnGID := args[1]

	if cfg.DryRun {
		out := output.NewJSON(os.Stdout)
		return out.Print(map[string]any{
			"dry_run":    true,
			"task_gid":   taskGID,
			"depends_on": dependsOnGID,
		})
	}

	client := newClient(cfg)
	if err := client.AddDependency(context.Background(), taskGID, dependsOnGID); err != nil {
		return err
	}

	out := output.NewJSON(os.Stdout)
	return out.Print(map[string]any{
		"task_gid":   taskGID,
		"depends_on": dependsOnGID,
		"created":    true,
	})
}

func runTaskDepList(_ *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	if err := requireAuth(cfg); err != nil {
		return err
	}

	taskGID := args[0]
	client := newClient(cfg)

	dependencies, err := client.ListDependencies(context.Background(), taskGID)
	if err != nil {
		return err
	}

	dependents, err := client.ListDependents(context.Background(), taskGID)
	if err != nil {
		return err
	}

	out := output.NewJSON(os.Stdout)
	return out.Print(map[string]any{
		"task_gid":   taskGID,
		"depends_on": dependencies,
		"dependents": dependents,
	})
}

func runTaskDepRm(_ *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	if err := requireAuth(cfg); err != nil {
		return err
	}

	taskGID := args[0]
	dependsOnGID := args[1]

	if cfg.DryRun {
		out := output.NewJSON(os.Stdout)
		return out.Print(map[string]any{
			"dry_run":      true,
			"task_gid":     taskGID,
			"removed_from": dependsOnGID,
		})
	}

	client := newClient(cfg)
	if err := client.RemoveDependency(context.Background(), taskGID, dependsOnGID); err != nil {
		return err
	}

	out := output.NewJSON(os.Stdout)
	return out.Print(map[string]any{
		"task_gid":     taskGID,
		"removed_from": dependsOnGID,
		"removed":      true,
	})
}
