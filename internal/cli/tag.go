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
	Long:  "List and get tags.",
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

var (
	tagListLimit  int
	tagListOffset string
)

func init() {
	rootCmd.AddCommand(tagCmd)
	tagCmd.AddCommand(tagListCmd)
	tagCmd.AddCommand(tagGetCmd)

	tagListCmd.Flags().IntVar(&tagListLimit, "limit", 50, "Max results to return")
	tagListCmd.Flags().StringVar(&tagListOffset, "offset", "", "Pagination offset")
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
