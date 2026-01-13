package cli

import (
	"context"
	"os"

	"github.com/spf13/cobra"

	"github.com/whoaa512/asana-cli/internal/errors"
	"github.com/whoaa512/asana-cli/internal/output"
)

var taskMoveCmd = &cobra.Command{
	Use:   "move <task-gid>",
	Short: "Move task to a section",
	Args:  cobra.ExactArgs(1),
	RunE:  runTaskMove,
}

var taskStartCmd = &cobra.Command{
	Use:   "start <task-gid>",
	Short: "Move task to in_progress section",
	Args:  cobra.ExactArgs(1),
	RunE:  runTaskStart,
}

var taskBlockCmd = &cobra.Command{
	Use:   "block <task-gid>",
	Short: "Move task to blocked section",
	Args:  cobra.ExactArgs(1),
	RunE:  runTaskBlock,
}

var (
	taskMoveSection string
)

func init() {
	taskCmd.AddCommand(taskMoveCmd)
	taskCmd.AddCommand(taskStartCmd)
	taskCmd.AddCommand(taskBlockCmd)

	taskMoveCmd.Flags().StringVar(&taskMoveSection, "section", "", "Section GID (required)")
	if err := taskMoveCmd.MarkFlagRequired("section"); err != nil {
		panic(err)
	}
}

func runTaskMove(_ *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	if err := requireAuth(cfg); err != nil {
		return err
	}

	taskGID := args[0]

	if cfg.DryRun {
		out := output.NewJSON(os.Stdout)
		return out.Print(map[string]any{"dry_run": true, "task": taskGID, "section": taskMoveSection})
	}

	client := newClient(cfg)
	if err := client.AddTaskToSection(context.Background(), taskMoveSection, taskGID); err != nil {
		return err
	}

	out := output.NewJSON(os.Stdout)
	return out.Print(map[string]any{"success": true, "task": taskGID, "section": taskMoveSection})
}

func runTaskStart(_ *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	if err := requireAuth(cfg); err != nil {
		return err
	}

	if cfg.Sections == nil || cfg.Sections["in_progress"] == "" {
		return errors.NewGeneralError("in_progress section not configured in .asana.json", nil)
	}

	taskGID := args[0]
	sectionGID := cfg.Sections["in_progress"]

	if cfg.DryRun {
		out := output.NewJSON(os.Stdout)
		return out.Print(map[string]any{"dry_run": true, "task": taskGID, "section": sectionGID, "section_name": "in_progress"})
	}

	client := newClient(cfg)
	if err := client.AddTaskToSection(context.Background(), sectionGID, taskGID); err != nil {
		return err
	}

	out := output.NewJSON(os.Stdout)
	return out.Print(map[string]any{"success": true, "task": taskGID, "section": sectionGID, "section_name": "in_progress"})
}

func runTaskBlock(_ *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	if err := requireAuth(cfg); err != nil {
		return err
	}

	if cfg.Sections == nil || cfg.Sections["blocked"] == "" {
		return errors.NewGeneralError("blocked section not configured in .asana.json", nil)
	}

	taskGID := args[0]
	sectionGID := cfg.Sections["blocked"]

	if cfg.DryRun {
		out := output.NewJSON(os.Stdout)
		return out.Print(map[string]any{"dry_run": true, "task": taskGID, "section": sectionGID, "section_name": "blocked"})
	}

	client := newClient(cfg)
	if err := client.AddTaskToSection(context.Background(), sectionGID, taskGID); err != nil {
		return err
	}

	out := output.NewJSON(os.Stdout)
	return out.Print(map[string]any{"success": true, "task": taskGID, "section": sectionGID, "section_name": "blocked"})
}
