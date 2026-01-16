package cli

import (
	"context"
	"os"

	"github.com/spf13/cobra"

	"github.com/whoaa512/asana-cli/internal/api"
	"github.com/whoaa512/asana-cli/internal/errors"
	"github.com/whoaa512/asana-cli/internal/output"
)

var meCmd = &cobra.Command{
	Use:   "me",
	Short: "Show current user info",
	Long:  "Display information about the authenticated user. Useful for verifying auth is working.",
	RunE:  runMe,
}

var meTeamsCmd = &cobra.Command{
	Use:   "teams",
	Short: "List my teams",
	Long:  "List teams the current user is a member of.",
	RunE:  runMeTeams,
}

var meProjectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "List workspace projects",
	Long:  "List projects in my workspace. Asana doesn't have user-specific project filtering; this returns all visible workspace projects.",
	RunE:  runMeProjects,
}

var meTasksCmd = &cobra.Command{
	Use:   "tasks",
	Short: "List my tasks",
	Long:  "List tasks assigned to the current user.",
	RunE:  runMeTasks,
}

var (
	meTeamsLimit     int
	meTeamsOffset    string
	meProjectsLimit  int
	meProjectsOffset string
	meTasksLimit     int
	meTasksOffset    string
	meTasksCompleted bool
)

func init() {
	rootCmd.AddCommand(meCmd)
	meCmd.AddCommand(meTeamsCmd)
	meCmd.AddCommand(meProjectsCmd)
	meCmd.AddCommand(meTasksCmd)

	meTeamsCmd.Flags().IntVar(&meTeamsLimit, "limit", 50, "Max results to return")
	meTeamsCmd.Flags().StringVar(&meTeamsOffset, "offset", "", "Pagination offset")

	meProjectsCmd.Flags().IntVar(&meProjectsLimit, "limit", 50, "Max results to return")
	meProjectsCmd.Flags().StringVar(&meProjectsOffset, "offset", "", "Pagination offset")

	meTasksCmd.Flags().IntVar(&meTasksLimit, "limit", 50, "Max results to return")
	meTasksCmd.Flags().StringVar(&meTasksOffset, "offset", "", "Pagination offset")
	meTasksCmd.Flags().BoolVar(&meTasksCompleted, "completed", false, "Show completed instead of incomplete tasks")
}

func runMe(_ *cobra.Command, _ []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	if err := requireAuth(cfg); err != nil {
		return err
	}

	client := newClient(cfg)
	user, err := client.GetMe(context.Background())
	if err != nil {
		return err
	}

	out := output.NewJSON(os.Stdout)
	return out.Print(user)
}

func runMeTeams(_ *cobra.Command, _ []string) error {
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

	opts := api.UserTeamListOptions{
		UserGID:      "me",
		Organization: cfg.Workspace,
		Limit:        meTeamsLimit,
		Offset:       meTeamsOffset,
	}

	client := newClient(cfg)
	result, err := client.ListUserTeams(context.Background(), opts)
	if err != nil {
		return err
	}

	out := output.NewJSON(os.Stdout)
	return out.Print(result)
}

func runMeProjects(_ *cobra.Command, _ []string) error {
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

	opts := api.UserProjectListOptions{
		Workspace: cfg.Workspace,
		Limit:     meProjectsLimit,
		Offset:    meProjectsOffset,
	}

	client := newClient(cfg)
	result, err := client.ListUserProjects(context.Background(), opts)
	if err != nil {
		return err
	}

	out := output.NewJSON(os.Stdout)
	return out.Print(result)
}

func runMeTasks(_ *cobra.Command, _ []string) error {
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

	opts := api.TaskListOptions{
		Workspace: cfg.Workspace,
		Assignee:  "me",
		Limit:     meTasksLimit,
		Offset:    meTasksOffset,
		Completed: &meTasksCompleted,
	}

	client := newClient(cfg)
	result, err := client.ListTasks(context.Background(), opts)
	if err != nil {
		return err
	}

	out := newOutput()
	return out.PrintTaskList(result)
}
