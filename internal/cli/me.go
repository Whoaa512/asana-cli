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

var (
	meTeamsLimit  int
	meTeamsOffset string
)

func init() {
	rootCmd.AddCommand(meCmd)
	meCmd.AddCommand(meTeamsCmd)

	meTeamsCmd.Flags().IntVar(&meTeamsLimit, "limit", 50, "Max results to return")
	meTeamsCmd.Flags().StringVar(&meTeamsOffset, "offset", "", "Pagination offset")
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
