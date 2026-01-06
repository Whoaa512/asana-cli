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

var sectionCmd = &cobra.Command{
	Use:   "section",
	Short: "Manage sections",
	Long:  "List, get, create sections and add tasks to sections.",
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

var sectionCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a section",
	RunE:  runSectionCreate,
}

var sectionAddTaskCmd = &cobra.Command{
	Use:   "add-task <section-gid> <task-gid>",
	Short: "Add a task to a section",
	Args:  cobra.ExactArgs(2),
	RunE:  runSectionAddTask,
}

var (
	sectionListProject   string
	sectionListLimit     int
	sectionListOffset    string
	sectionCreateProject string
	sectionCreateName    string
)

func init() {
	rootCmd.AddCommand(sectionCmd)
	sectionCmd.AddCommand(sectionListCmd)
	sectionCmd.AddCommand(sectionGetCmd)
	sectionCmd.AddCommand(sectionCreateCmd)
	sectionCmd.AddCommand(sectionAddTaskCmd)

	sectionListCmd.Flags().StringVar(&sectionListProject, "project", "", "Project GID")
	sectionListCmd.Flags().IntVar(&sectionListLimit, "limit", 50, "Max results to return")
	sectionListCmd.Flags().StringVar(&sectionListOffset, "offset", "", "Pagination offset")

	sectionCreateCmd.Flags().StringVar(&sectionCreateProject, "project", "", "Project GID")
	sectionCreateCmd.Flags().StringVar(&sectionCreateName, "name", "", "Section name (required)")
	if err := sectionCreateCmd.MarkFlagRequired("name"); err != nil {
		panic(err)
	}
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

func runSectionCreate(_ *cobra.Command, _ []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	if err := requireAuth(cfg); err != nil {
		return err
	}

	project := sectionCreateProject
	if project == "" && cfg.Project != "" {
		project = cfg.Project
	}
	if project == "" {
		return errors.NewGeneralError("no project specified (use --project or set in .asana.json)", nil)
	}

	req := models.SectionCreateRequest{
		Name: sectionCreateName,
	}

	if cfg.DryRun {
		out := output.NewJSON(os.Stdout)
		return out.Print(map[string]any{"dry_run": true, "project": project, "request": req})
	}

	client := newClient(cfg)
	section, err := client.CreateSection(context.Background(), project, req)
	if err != nil {
		return err
	}

	out := output.NewJSON(os.Stdout)
	return out.Print(section)
}

func runSectionAddTask(_ *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	if err := requireAuth(cfg); err != nil {
		return err
	}

	sectionGID := args[0]
	taskGID := args[1]

	if cfg.DryRun {
		out := output.NewJSON(os.Stdout)
		return out.Print(map[string]any{"dry_run": true, "section": sectionGID, "task": taskGID})
	}

	client := newClient(cfg)
	if err := client.AddTaskToSection(context.Background(), sectionGID, taskGID); err != nil {
		return err
	}

	out := output.NewJSON(os.Stdout)
	return out.Print(map[string]any{"success": true, "section": sectionGID, "task": taskGID})
}
