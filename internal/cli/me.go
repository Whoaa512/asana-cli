package cli

import (
	"context"
	"os"

	"github.com/spf13/cobra"

	"github.com/whoaa512/asana-cli/internal/output"
)

var meCmd = &cobra.Command{
	Use:   "me",
	Short: "Show current user info",
	Long:  "Display information about the authenticated user. Useful for verifying auth is working.",
	RunE:  runMe,
}

func init() {
	rootCmd.AddCommand(meCmd)
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
