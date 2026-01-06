package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "asana",
	Short: "CLI for interacting with Asana",
	Long:  "A CLI tool for managing Asana tasks, designed for AI agents with JSON-only output.",
}

func Execute() int {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	return 0
}
