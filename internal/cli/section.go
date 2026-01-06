package cli

import (
	"context"
	"os"

	"github.com/spf13/cobra"

	"github.com/whoaa512/asana-cli/internal/api"
	"github.com/whoaa512/asana-cli/internal/errors"
	"github.com/whoaa512/asana-cli/internal/output"
)

var sectionCmd = &cobra.Command{
	Use:   "section",
	Short: "Manage sections",
	Long:  "List and get section details.",
}

var sectionListCmd = &cobra.Command{
	Use:   "list",
	Short: "List sections",
	Long:  "List sections in a project.",
	RunE:  runSectionList,
}

var sectionGetCmd = &cobra.Command{
	Use:   "get <gid>",
	Short: "Get section details",
	Args:  cobra.ExactArgs(1),
	RunE:  runSectionGet,
}

var (
	sectionListProject string
	sectionListLimit   int
	sectionListOffset  string
)

func init() {
	rootCmd.AddCommand(sectionCmd)
	sectionCmd.AddCommand(sectionListCmd)
	sectionCmd.AddCommand(sectionGetCmd)

	sectionListCmd.Flags().StringVar(&sectionListProject, "project", "", "Project GID (required)")
	sectionListCmd.Flags().IntVar(&sectionListLimit, "limit", 50, "Max results to return")
	sectionListCmd.Flags().StringVar(&sectionListOffset, "offset", "", "Pagination offset")
}

func runSectionList(_ *cobra.Command, _ []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	if err := requireAuth(cfg); err != nil {
		return err
	}

	project := sectionListProject
	if project == "" && cfg.Project != "" {
		project = cfg.Project
	}
	if project == "" {
		return errors.NewGeneralError("no project specified (use --project or set in .asana.json)", nil)
	}

	opts := api.SectionListOptions{
		Project: project,
		Limit:   sectionListLimit,
		Offset:  sectionListOffset,
	}

	client := newClient(cfg)
	result, err := client.ListSections(context.Background(), opts)
	if err != nil {
		return err
	}

	out := output.NewJSON(os.Stdout)
	return out.Print(result)
}

func runSectionGet(_ *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	if err := requireAuth(cfg); err != nil {
		return err
	}

	client := newClient(cfg)
	section, err := client.GetSection(context.Background(), args[0])
	if err != nil {
		return err
	}

	out := output.NewJSON(os.Stdout)
	return out.Print(section)
}
