package cli

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/whoaa512/asana-cli/internal/output"
)

var cliVersion = "dev"

func SetVersion(v string) {
	cliVersion = v
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long:  "Print the version of the asana-cli tool.",
	RunE:  runVersion,
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func runVersion(_ *cobra.Command, _ []string) error {
	out := output.NewJSON(os.Stdout)
	return out.Print(map[string]string{
		"version": cliVersion,
	})
}
