package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/whoaa512/asana-cli/internal/config"
	"github.com/whoaa512/asana-cli/internal/errors"
	"github.com/whoaa512/asana-cli/internal/models"
	"github.com/whoaa512/asana-cli/internal/output"
)

var workspaceCmd = &cobra.Command{
	Use:   "workspace",
	Short: "Manage workspaces",
	Long:  "List, get, and set the active workspace.",
}

var workspaceListCmd = &cobra.Command{
	Use:   "list",
	Short: "List workspaces",
	RunE:  runWorkspaceList,
}

var workspaceGetCmd = &cobra.Command{
	Use:   "get <gid>",
	Short: "Get workspace details",
	Args:  cobra.ExactArgs(1),
	RunE:  runWorkspaceGet,
}

var workspaceUseCmd = &cobra.Command{
	Use:   "use <gid>",
	Short: "Set active workspace",
	Long:  "Set the active workspace. Without --global, writes to local .asana.json. With --global, writes to ~/.config/asana-cli/config.json.",
	Args:  cobra.ExactArgs(1),
	RunE:  runWorkspaceUse,
}

var (
	workspaceListLimit int
	workspaceUseGlobal bool
)

func init() {
	rootCmd.AddCommand(workspaceCmd)
	workspaceCmd.AddCommand(workspaceListCmd)
	workspaceCmd.AddCommand(workspaceGetCmd)
	workspaceCmd.AddCommand(workspaceUseCmd)

	workspaceListCmd.Flags().IntVar(&workspaceListLimit, "limit", 50, "Max results to return")
	workspaceUseCmd.Flags().BoolVar(&workspaceUseGlobal, "global", false, "Write to global config instead of local .asana.json")
}

func runWorkspaceList(_ *cobra.Command, _ []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	if err := requireAuth(cfg); err != nil {
		return err
	}

	client := newClient(cfg)
	result, err := client.ListWorkspaces(context.Background(), workspaceListLimit)
	if err != nil {
		return err
	}

	out := output.NewJSON(os.Stdout)
	return out.Print(result)
}

func runWorkspaceGet(_ *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	if err := requireAuth(cfg); err != nil {
		return err
	}

	client := newClient(cfg)
	workspace, err := client.GetWorkspace(context.Background(), args[0])
	if err != nil {
		return err
	}

	out := output.NewJSON(os.Stdout)
	return out.Print(workspace)
}

func runWorkspaceUse(_ *cobra.Command, args []string) error {
	gid := args[0]

	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	if err := requireAuth(cfg); err != nil {
		return err
	}

	client := newClient(cfg)
	workspace, err := client.GetWorkspace(context.Background(), gid)
	if err != nil {
		return err
	}

	if cfg.DryRun {
		return printWorkspaceUseDryRun(cfg, gid, workspace)
	}

	if workspaceUseGlobal {
		if err := setGlobalWorkspace(cfg.ConfigPath, gid); err != nil {
			return err
		}
	} else {
		if err := setLocalWorkspace(gid); err != nil {
			return err
		}
	}

	out := output.NewJSON(os.Stdout)
	return out.Print(workspace)
}

func setGlobalWorkspace(configPath, gid string) error {
	var data map[string]any
	content, err := os.ReadFile(configPath)
	if err != nil && !os.IsNotExist(err) {
		return errors.NewGeneralError("failed to read config", err)
	}
	if len(content) > 0 {
		if err := json.Unmarshal(content, &data); err != nil {
			return errors.NewGeneralError("failed to parse config", err)
		}
	}
	if data == nil {
		data = make(map[string]any)
	}

	data["default_workspace"] = gid

	newContent, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return errors.NewGeneralError("failed to encode config", err)
	}

	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return errors.NewGeneralError("failed to create config directory", err)
	}

	if err := os.WriteFile(configPath, newContent, 0644); err != nil {
		return errors.NewGeneralError("failed to write config", err)
	}

	return nil
}

func setLocalWorkspace(gid string) error {
	ctx, err := config.LoadLocalContext()
	if err != nil {
		return errors.NewGeneralError("failed to load local context", err)
	}

	ctx.Workspace = gid

	dir, err := config.FindContextFileDir()
	if err != nil {
		return errors.NewGeneralError("failed to find context directory", err)
	}

	if err := ctx.Save(dir); err != nil {
		return errors.NewGeneralError("failed to save local context", err)
	}

	return nil
}

func printWorkspaceUseDryRun(cfg *config.Config, gid string, ws *models.Workspace) error {
	fmt.Println("[dry-run] Would set workspace to:", ws.Name, "("+gid+")")

	if workspaceUseGlobal {
		fmt.Println("[dry-run] Target file:", cfg.ConfigPath)
		fmt.Println("[dry-run] Would write: default_workspace =", gid)
	} else {
		dir, err := config.FindContextFileDir()
		if err != nil {
			dir = "."
		}
		fmt.Println("[dry-run] Target file:", filepath.Join(dir, ".asana.json"))
		fmt.Println("[dry-run] Would write: workspace =", gid)
	}

	return nil
}
