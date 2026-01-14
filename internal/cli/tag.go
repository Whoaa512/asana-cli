package cli

import (
	"context"
	"os"

	"github.com/spf13/cobra"

	"github.com/whoaa512/asana-cli/internal/api"
	"github.com/whoaa512/asana-cli/internal/errors"
	"github.com/whoaa512/asana-cli/internal/output"
)

var tagCmd = &cobra.Command{
	Use:   "tag",
	Short: "Manage tags",
	Long:  "List, get, and create tags.",
}

var tagListCmd = &cobra.Command{
	Use:   "list",
	Short: "List tags",
	Long:  "List tags in a workspace.",
	RunE:  runTagList,
}

var tagGetCmd = &cobra.Command{
	Use:   "get <gid>",
	Short: "Get tag details",
	Args:  cobra.ExactArgs(1),
	RunE:  runTagGet,
}

var tagCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a tag",
	Long: `Create a new tag in a workspace.

Requires --name. Workspace can be specified via --workspace flag or from context.`,
	Example: `  # Create a tag using workspace from context
  asana tag create --name "urgent"

  # Create with explicit workspace
  asana tag create --name "bug" --workspace 1234567890

  # Create with color
  asana tag create --name "priority" --color dark-red`,
	RunE: runTagCreate,
}

var (
	tagListLimit  int
	tagListOffset string

	tagCreateName      string
	tagCreateWorkspace string
	tagCreateColor     string
)

func init() {
	rootCmd.AddCommand(tagCmd)
	tagCmd.AddCommand(tagListCmd)
	tagCmd.AddCommand(tagGetCmd)
	tagCmd.AddCommand(tagCreateCmd)

	tagListCmd.Flags().IntVar(&tagListLimit, "limit", 50, "Max results to return")
	tagListCmd.Flags().StringVar(&tagListOffset, "offset", "", "Pagination offset")

	tagCreateCmd.Flags().StringVar(&tagCreateName, "name", "", "Tag name (required)")
	tagCreateCmd.Flags().StringVar(&tagCreateWorkspace, "workspace", "", "Workspace GID")
	tagCreateCmd.Flags().StringVar(&tagCreateColor, "color", "", "Tag color (e.g., dark-red, light-blue)")
	_ = tagCreateCmd.MarkFlagRequired("name")
}

func runTagList(_ *cobra.Command, _ []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	if err := requireAuth(cfg); err != nil {
		return err
	}

	if cfg.Workspace == "" {
		return errors.NewGeneralError("no workspace specified", nil)
	}

	opts := api.TagListOptions{
		Workspace: cfg.Workspace,
		Limit:     tagListLimit,
		Offset:    tagListOffset,
	}

	client := newClient(cfg)
	result, err := client.ListTags(context.Background(), opts)
	if err != nil {
		return err
	}

	out := output.NewJSON(os.Stdout)
	return out.Print(result)
}

func runTagGet(_ *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	if err := requireAuth(cfg); err != nil {
		return err
	}

	client := newClient(cfg)
	tag, err := client.GetTag(context.Background(), args[0])
	if err != nil {
		return err
	}

	out := output.NewJSON(os.Stdout)
	return out.Print(tag)
}

func runTagCreate(_ *cobra.Command, _ []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	if err := requireAuth(cfg); err != nil {
		return err
	}

	workspace := tagCreateWorkspace
	if workspace == "" {
		workspace = cfg.Workspace
	}
	if workspace == "" {
		return errors.NewGeneralError("no workspace specified", nil)
	}

	req := api.TagCreateRequest{
		Name:      tagCreateName,
		Workspace: workspace,
		Color:     tagCreateColor,
	}

	if cfg.DryRun {
		out := output.NewJSON(os.Stdout)
		return out.Print(map[string]any{"dry_run": true, "request": req})
	}

	client := newClient(cfg)
	tag, err := client.CreateTag(context.Background(), req)
	if err != nil {
		return err
	}

	out := output.NewJSON(os.Stdout)
	return out.Print(tag)
}
