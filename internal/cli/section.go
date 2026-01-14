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

var sectionUpdateCmd = &cobra.Command{
	Use:   "update <gid>",
	Short: "Update a section",
	Args:  cobra.ExactArgs(1),
	RunE:  runSectionUpdate,
}

var sectionDeleteCmd = &cobra.Command{
	Use:   "delete <gid>",
	Short: "Delete a section",
	Args:  cobra.ExactArgs(1),
	RunE:  runSectionDelete,
}

var sectionInsertCmd = &cobra.Command{
	Use:   "insert <section-gid>",
	Short: "Reorder a section within a project",
	Args:  cobra.ExactArgs(1),
	RunE:  runSectionInsert,
}

var (
	sectionListProject   string
	sectionListLimit     int
	sectionListOffset    string
	sectionCreateProject string
	sectionCreateName    string
	sectionUpdateName    string
	sectionInsertProject string
	sectionInsertBefore  string
	sectionInsertAfter   string
)

func init() {
	rootCmd.AddCommand(sectionCmd)
	sectionCmd.AddCommand(sectionListCmd)
	sectionCmd.AddCommand(sectionGetCmd)
	sectionCmd.AddCommand(sectionCreateCmd)
	sectionCmd.AddCommand(sectionUpdateCmd)
	sectionCmd.AddCommand(sectionDeleteCmd)
	sectionCmd.AddCommand(sectionInsertCmd)
	sectionCmd.AddCommand(sectionAddTaskCmd)

	sectionListCmd.Flags().StringVar(&sectionListProject, "project", "", "Project GID")
	sectionListCmd.Flags().IntVar(&sectionListLimit, "limit", 50, "Max results to return")
	sectionListCmd.Flags().StringVar(&sectionListOffset, "offset", "", "Pagination offset")

	sectionCreateCmd.Flags().StringVar(&sectionCreateProject, "project", "", "Project GID")
	sectionCreateCmd.Flags().StringVar(&sectionCreateName, "name", "", "Section name (required)")
	if err := sectionCreateCmd.MarkFlagRequired("name"); err != nil {
		panic(err)
	}

	sectionUpdateCmd.Flags().StringVar(&sectionUpdateName, "name", "", "New section name (required)")
	if err := sectionUpdateCmd.MarkFlagRequired("name"); err != nil {
		panic(err)
	}

	sectionInsertCmd.Flags().StringVar(&sectionInsertProject, "project", "", "Project GID")
	sectionInsertCmd.Flags().StringVar(&sectionInsertBefore, "before", "", "Insert before this section GID")
	sectionInsertCmd.Flags().StringVar(&sectionInsertAfter, "after", "", "Insert after this section GID")
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

func runSectionUpdate(_ *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	if err := requireAuth(cfg); err != nil {
		return err
	}

	req := models.SectionUpdateRequest{
		Name: sectionUpdateName,
	}

	if cfg.DryRun {
		out := output.NewJSON(os.Stdout)
		return out.Print(map[string]any{"dry_run": true, "gid": args[0], "request": req})
	}

	client := newClient(cfg)
	section, err := client.UpdateSection(context.Background(), args[0], req)
	if err != nil {
		return err
	}

	out := output.NewJSON(os.Stdout)
	return out.Print(section)
}

func runSectionDelete(_ *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	if err := requireAuth(cfg); err != nil {
		return err
	}

	if cfg.DryRun {
		out := output.NewJSON(os.Stdout)
		return out.Print(map[string]any{"dry_run": true, "gid": args[0], "action": "delete"})
	}

	client := newClient(cfg)
	if err := client.DeleteSection(context.Background(), args[0]); err != nil {
		return err
	}

	out := output.NewJSON(os.Stdout)
	return out.Print(map[string]any{"deleted": true, "gid": args[0]})
}

func runSectionInsert(_ *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	if err := requireAuth(cfg); err != nil {
		return err
	}

	project := sectionInsertProject
	if project == "" && cfg.Project != "" {
		project = cfg.Project
	}
	if project == "" {
		return errors.NewGeneralError("no project specified (use --project or set in .asana.json)", nil)
	}

	if sectionInsertBefore == "" && sectionInsertAfter == "" {
		return errors.NewGeneralError("must specify either --before or --after", nil)
	}
	if sectionInsertBefore != "" && sectionInsertAfter != "" {
		return errors.NewGeneralError("cannot specify both --before and --after", nil)
	}

	req := models.SectionInsertRequest{
		Section: args[0],
	}
	if sectionInsertBefore != "" {
		req.BeforeSection = &sectionInsertBefore
	}
	if sectionInsertAfter != "" {
		req.AfterSection = &sectionInsertAfter
	}

	if cfg.DryRun {
		out := output.NewJSON(os.Stdout)
		return out.Print(map[string]any{"dry_run": true, "project": project, "request": req})
	}

	client := newClient(cfg)
	if err := client.InsertSection(context.Background(), project, req); err != nil {
		return err
	}

	out := output.NewJSON(os.Stdout)
	return out.Print(map[string]any{"success": true, "project": project, "section": args[0]})
}
