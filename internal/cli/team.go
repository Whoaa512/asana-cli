package cli

import (
	"context"
	"os"

	"github.com/spf13/cobra"

	"github.com/whoaa512/asana-cli/internal/api"
	"github.com/whoaa512/asana-cli/internal/errors"
	"github.com/whoaa512/asana-cli/internal/output"
)

var teamCmd = &cobra.Command{
	Use:   "team",
	Short: "Manage teams",
	Long:  "List teams in an organization.",
}

var teamListCmd = &cobra.Command{
	Use:   "list",
	Short: "List teams",
	Long:  "List teams in an organization workspace.",
	RunE:  runTeamList,
}

var teamMeCmd = &cobra.Command{
	Use:   "me",
	Short: "List my teams",
	Long:  "List teams the current user is a member of.",
	RunE:  runTeamMe,
}

var (
	teamListLimit  int
	teamListOffset string
	teamMeLimit    int
	teamMeOffset   string
)

func init() {
	rootCmd.AddCommand(teamCmd)
	teamCmd.AddCommand(teamListCmd)
	teamCmd.AddCommand(teamMeCmd)

	teamListCmd.Flags().IntVar(&teamListLimit, "limit", 50, "Max results to return")
	teamListCmd.Flags().StringVar(&teamListOffset, "offset", "", "Pagination offset")

	teamMeCmd.Flags().IntVar(&teamMeLimit, "limit", 50, "Max results to return")
	teamMeCmd.Flags().StringVar(&teamMeOffset, "offset", "", "Pagination offset")
}

func runTeamList(_ *cobra.Command, _ []string) error {
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

	opts := api.TeamListOptions{
		Organization: cfg.Workspace,
		Limit:        teamListLimit,
		Offset:       teamListOffset,
	}

	client := newClient(cfg)
	result, err := client.ListTeams(context.Background(), opts)
	if err != nil {
		return err
	}

	out := output.NewJSON(os.Stdout)
	return out.Print(result)
}

func runTeamMe(_ *cobra.Command, _ []string) error {
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
		Limit:        teamMeLimit,
		Offset:       teamMeOffset,
	}

	client := newClient(cfg)
	result, err := client.ListUserTeams(context.Background(), opts)
	if err != nil {
		return err
	}

	out := output.NewJSON(os.Stdout)
	return out.Print(result)
}
