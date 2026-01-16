package cli

import (
	"context"
	"os"

	"github.com/spf13/cobra"

	"github.com/whoaa512/asana-cli/internal/config"
	"github.com/whoaa512/asana-cli/internal/errors"
	"github.com/whoaa512/asana-cli/internal/output"
)

var ctxCmd = &cobra.Command{
	Use:   "ctx",
	Short: "Manage local context",
	Long:  "View and modify the local .asana.json context file.",
}

var ctxShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current context",
	RunE:  runCtxShow,
}

var ctxTaskCmd = &cobra.Command{
	Use:   "task [<task>]",
	Short: "Get or set context task",
	Long:  "Without args, shows current task. With GID or name, sets task. With --clear, removes task. Uses fuzzy matching for names.",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runCtxTask,
}

var ctxProjectCmd = &cobra.Command{
	Use:   "project [<gid>]",
	Short: "Get or set context project",
	Long:  "Without args, shows current project. With gid, sets project. With --clear, removes project.",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runCtxProject,
}

var ctxClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear all context",
	RunE:  runCtxClear,
}

var (
	ctxTaskClear    bool
	ctxTaskPick     bool
	ctxProjectClear bool
)

func init() {
	rootCmd.AddCommand(ctxCmd)
	ctxCmd.AddCommand(ctxShowCmd)
	ctxCmd.AddCommand(ctxTaskCmd)
	ctxCmd.AddCommand(ctxProjectCmd)
	ctxCmd.AddCommand(ctxClearCmd)

	ctxTaskCmd.Flags().BoolVar(&ctxTaskClear, "clear", false, "Clear the task from context")
	ctxTaskCmd.Flags().BoolVar(&ctxTaskPick, "pick", false, "Show interactive picker if multiple matches")
	ctxProjectCmd.Flags().BoolVar(&ctxProjectClear, "clear", false, "Clear the project from context")
}

func runCtxShow(_ *cobra.Command, _ []string) error {
	ctx, err := config.LoadLocalContext()
	if err != nil {
		return errors.NewGeneralError("failed to load context", err)
	}

	result := map[string]any{
		"workspace": ctx.Workspace,
		"project":   ctx.Project,
		"task":      ctx.Task,
		"path":      ctx.Path(),
	}

	out := output.NewJSON(os.Stdout)
	return out.Print(result)
}

func runCtxTask(_ *cobra.Command, args []string) error {
	localCtx, err := config.LoadLocalContext()
	if err != nil {
		return errors.NewGeneralError("failed to load context", err)
	}

	if ctxTaskClear {
		localCtx.Task = ""
		return saveContext(localCtx)
	}

	if len(args) == 0 {
		result := map[string]string{"task": localCtx.Task}
		out := output.NewJSON(os.Stdout)
		return out.Print(result)
	}

	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	if err := requireAuth(cfg); err != nil {
		return err
	}

	client := newClient(cfg)
	ctx := context.Background()

	taskGID, err := resolveTaskGID(ctx, cfg, client, args[0], ctxTaskPick)
	if err != nil {
		return err
	}

	localCtx.Task = taskGID
	return saveContext(localCtx)
}

func runCtxProject(_ *cobra.Command, args []string) error {
	ctx, err := config.LoadLocalContext()
	if err != nil {
		return errors.NewGeneralError("failed to load context", err)
	}

	if ctxProjectClear {
		ctx.Project = ""
		return saveContext(ctx)
	}

	if len(args) == 0 {
		result := map[string]string{"project": ctx.Project}
		out := output.NewJSON(os.Stdout)
		return out.Print(result)
	}

	ctx.Project = args[0]
	return saveContext(ctx)
}

func runCtxClear(_ *cobra.Command, _ []string) error {
	ctx, err := config.LoadLocalContext()
	if err != nil {
		return errors.NewGeneralError("failed to load context", err)
	}

	ctx.Workspace = ""
	ctx.Project = ""
	ctx.Task = ""

	return saveContext(ctx)
}

func saveContext(ctx *config.LocalContext) error {
	dir, err := config.FindContextFileDir()
	if err != nil {
		return errors.NewGeneralError("failed to find context directory", err)
	}

	if err := ctx.Save(dir); err != nil {
		return errors.NewGeneralError("failed to save context", err)
	}

	out := output.NewJSON(os.Stdout)
	return out.Print(map[string]any{
		"workspace": ctx.Workspace,
		"project":   ctx.Project,
		"task":      ctx.Task,
		"path":      ctx.Path(),
	})
}
