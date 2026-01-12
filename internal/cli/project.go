package cli

import (
	"context"
	"os"

	"github.com/spf13/cobra"

	"github.com/whoaa512/asana-cli/internal/api"
	"github.com/whoaa512/asana-cli/internal/errors"
	"github.com/whoaa512/asana-cli/internal/models"
	"github.com/whoaa512/asana-cli/internal/output"
)

var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Manage projects",
	Long:  "List, get, and create projects.",
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

var projectCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a project",
	Long:  "Create a new project in Asana. Requires --name. Workspace is taken from config.",
	Example: `  # Create a project
  asana project create --name "Q1 Planning"

  # Create with notes and color
  asana project create --name "Q1 Planning" --notes "Planning for Q1 2024" --color "light-green"`,
	RunE: runProjectCreate,
}

var (
	projectListArchived bool
	projectListLimit    int
	projectListOffset   string

	projectCreateName  string
	projectCreateNotes string
	projectCreateColor string
	projectCreateTeam  string
)

func init() {
	rootCmd.AddCommand(projectCmd)
	projectCmd.AddCommand(projectListCmd)
	projectCmd.AddCommand(projectGetCmd)
	projectCmd.AddCommand(projectCreateCmd)

	projectListCmd.Flags().BoolVar(&projectListArchived, "archived", false, "Include archived projects")
	projectListCmd.Flags().IntVar(&projectListLimit, "limit", 50, "Max results to return")
	projectListCmd.Flags().StringVar(&projectListOffset, "offset", "", "Pagination offset")

	projectCreateCmd.Flags().StringVar(&projectCreateName, "name", "", "Project name (required)")
	projectCreateCmd.Flags().StringVar(&projectCreateNotes, "notes", "", "Project description")
	projectCreateCmd.Flags().StringVar(&projectCreateColor, "color", "", "Project color")
	projectCreateCmd.Flags().StringVar(&projectCreateTeam, "team", "", "Team GID (required for org workspaces)")
	_ = projectCreateCmd.MarkFlagRequired("name")
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

func runProjectCreate(_ *cobra.Command, _ []string) error {
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

	req := models.ProjectCreateRequest{
		Name:  projectCreateName,
		Notes: projectCreateNotes,
		Color: projectCreateColor,
	}
	team := projectCreateTeam
	if team == "" {
		team = cfg.Team
	}
	if team != "" {
		req.Team = team
	} else {
		req.Workspace = cfg.Workspace
	}

	if cfg.DryRun {
		out := output.NewJSON(os.Stdout)
		return out.Print(map[string]any{"dry_run": true, "request": req})
	}

	client := newClient(cfg)
	project, err := client.CreateProject(context.Background(), req)
	if err != nil {
		return err
	}

	out := output.NewJSON(os.Stdout)
	return out.Print(project)
}
