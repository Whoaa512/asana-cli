package cli

import (
	"context"
	"os"

	"github.com/spf13/cobra"

	"github.com/whoaa512/asana-cli/internal/api"
	"github.com/whoaa512/asana-cli/internal/errors"
	"github.com/whoaa512/asana-cli/internal/output"
)

var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Manage projects",
	Long:  "List and get project details.",
}

var projectListCmd = &cobra.Command{
	Use:   "list",
	Short: "List projects",
	Long:  "List projects in a workspace.",
	RunE:  runProjectList,
}

var projectGetCmd = &cobra.Command{
	Use:   "get <gid>",
	Short: "Get project details",
	Args:  cobra.ExactArgs(1),
	RunE:  runProjectGet,
}

var (
	projectListArchived bool
	projectListLimit    int
	projectListOffset   string
)

func init() {
	rootCmd.AddCommand(projectCmd)
	projectCmd.AddCommand(projectListCmd)
	projectCmd.AddCommand(projectGetCmd)

	projectListCmd.Flags().BoolVar(&projectListArchived, "archived", false, "Include archived projects")
	projectListCmd.Flags().IntVar(&projectListLimit, "limit", 50, "Max results to return")
	projectListCmd.Flags().StringVar(&projectListOffset, "offset", "", "Pagination offset")
}

func runProjectList(_ *cobra.Command, _ []string) error {
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

	opts := api.ProjectListOptions{
		Workspace: cfg.Workspace,
		Archived:  projectListArchived,
		Limit:     projectListLimit,
		Offset:    projectListOffset,
	}

	client := newClient(cfg)
	result, err := client.ListProjects(context.Background(), opts)
	if err != nil {
		return err
	}

	out := output.NewJSON(os.Stdout)
	return out.Print(result)
}

func runProjectGet(_ *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	if err := requireAuth(cfg); err != nil {
		return err
	}

	client := newClient(cfg)
	project, err := client.GetProject(context.Background(), args[0])
	if err != nil {
		return err
	}

	out := output.NewJSON(os.Stdout)
	return out.Print(project)
}
