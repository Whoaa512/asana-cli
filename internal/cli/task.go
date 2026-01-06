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

var taskCmd = &cobra.Command{
	Use:   "task",
	Short: "Manage tasks",
	Long:  "List, get, create, update, complete, and delete tasks.",
}

var taskListCmd = &cobra.Command{
	Use:   "list",
	Short: "List tasks",
	Long:  "List tasks. Requires --project or workspace from context.",
	RunE:  runTaskList,
}

var taskGetCmd = &cobra.Command{
	Use:   "get <gid>",
	Short: "Get task details",
	Args:  cobra.ExactArgs(1),
	RunE:  runTaskGet,
}

var taskCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a task",
	Long: `Create a new task in Asana.

Requires --name. Project can be specified via --project flag or from .asana.json context.
Use --parent to create a subtask instead.`,
	Example: `  # Create a task in a project
  asana task create --name "Fix login bug" --project 1234567890

  # Create with due date and assignee
  asana task create --name "Review PR" --due-on 2024-01-15 --assignee me

  # Create a subtask
  asana task create --name "Write tests" --parent 9876543210`,
	RunE: runTaskCreate,
}

var taskUpdateCmd = &cobra.Command{
	Use:   "update <gid>",
	Short: "Update a task",
	Args:  cobra.ExactArgs(1),
	RunE:  runTaskUpdate,
}

var taskCompleteCmd = &cobra.Command{
	Use:   "complete <gid>",
	Short: "Mark task as complete",
	Args:  cobra.ExactArgs(1),
	RunE:  runTaskComplete,
}

var taskDeleteCmd = &cobra.Command{
	Use:   "delete <gid>",
	Short: "Delete a task",
	Args:  cobra.ExactArgs(1),
	RunE:  runTaskDelete,
}

var taskAssignCmd = &cobra.Command{
	Use:   "assign <gid> <assignee>",
	Short: "Assign a task to a user",
	Args:  cobra.ExactArgs(2),
	RunE:  runTaskAssign,
}

var (
	taskListProject   string
	taskListAssignee  string
	taskListCompleted string
	taskListLimit     int
	taskListOffset    string

	taskCreateName     string
	taskCreateNotes    string
	taskCreateProject  string
	taskCreateAssignee string
	taskCreateDueOn    string
	taskCreateParent   string

	taskUpdateName     string
	taskUpdateNotes    string
	taskUpdateAssignee string
	taskUpdateDueOn    string
)

func init() {
	rootCmd.AddCommand(taskCmd)
	taskCmd.AddCommand(taskListCmd)
	taskCmd.AddCommand(taskGetCmd)
	taskCmd.AddCommand(taskCreateCmd)
	taskCmd.AddCommand(taskUpdateCmd)
	taskCmd.AddCommand(taskCompleteCmd)
	taskCmd.AddCommand(taskDeleteCmd)
	taskCmd.AddCommand(taskAssignCmd)

	taskListCmd.Flags().StringVar(&taskListProject, "project", "", "Filter by project GID")
	taskListCmd.Flags().StringVar(&taskListAssignee, "assignee", "", "Filter by assignee GID or 'me'")
	taskListCmd.Flags().StringVar(&taskListCompleted, "completed", "", "Filter by completed status (true/false)")
	taskListCmd.Flags().IntVar(&taskListLimit, "limit", 50, "Max results to return")
	taskListCmd.Flags().StringVar(&taskListOffset, "offset", "", "Pagination offset")

	taskCreateCmd.Flags().StringVar(&taskCreateName, "name", "", "Task name (required)")
	taskCreateCmd.Flags().StringVar(&taskCreateNotes, "notes", "", "Task notes/description")
	taskCreateCmd.Flags().StringVar(&taskCreateProject, "project", "", "Project GID")
	taskCreateCmd.Flags().StringVar(&taskCreateAssignee, "assignee", "", "Assignee GID or 'me'")
	taskCreateCmd.Flags().StringVar(&taskCreateDueOn, "due-on", "", "Due date (YYYY-MM-DD)")
	taskCreateCmd.Flags().StringVar(&taskCreateParent, "parent", "", "Parent task GID (for subtasks)")
	_ = taskCreateCmd.MarkFlagRequired("name")

	taskUpdateCmd.Flags().StringVar(&taskUpdateName, "name", "", "New task name")
	taskUpdateCmd.Flags().StringVar(&taskUpdateNotes, "notes", "", "New task notes")
	taskUpdateCmd.Flags().StringVar(&taskUpdateAssignee, "assignee", "", "New assignee GID or 'me'")
	taskUpdateCmd.Flags().StringVar(&taskUpdateDueOn, "due-on", "", "New due date (YYYY-MM-DD)")
}

func runTaskList(_ *cobra.Command, _ []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	if err := requireAuth(cfg); err != nil {
		return err
	}

	opts := api.TaskListOptions{
		Project:  taskListProject,
		Assignee: taskListAssignee,
		Limit:    taskListLimit,
		Offset:   taskListOffset,
	}

	if opts.Project == "" && cfg.Project != "" {
		opts.Project = cfg.Project
	}

	if opts.Project == "" {
		if cfg.Workspace == "" {
			return errors.NewGeneralError("no project or workspace specified", nil)
		}
		opts.Workspace = cfg.Workspace
	}

	if taskListCompleted != "" {
		completed := taskListCompleted == "true"
		opts.Completed = &completed
	}

	client := newClient(cfg)
	result, err := client.ListTasks(context.Background(), opts)
	if err != nil {
		return err
	}

	out := output.NewJSON(os.Stdout)
	return out.Print(result)
}

func runTaskGet(_ *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	if err := requireAuth(cfg); err != nil {
		return err
	}

	client := newClient(cfg)
	task, err := client.GetTask(context.Background(), args[0])
	if err != nil {
		return err
	}

	out := output.NewJSON(os.Stdout)
	return out.Print(task)
}

func runTaskCreate(_ *cobra.Command, _ []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	if err := requireAuth(cfg); err != nil {
		return err
	}

	req := models.TaskCreateRequest{
		Name:     taskCreateName,
		Notes:    taskCreateNotes,
		Assignee: taskCreateAssignee,
		DueOn:    taskCreateDueOn,
		Parent:   taskCreateParent,
	}

	project := taskCreateProject
	if project == "" && cfg.Project != "" {
		project = cfg.Project
	}
	if project != "" {
		req.Projects = []string{project}
	}

	if req.Parent == "" && req.Projects == nil {
		if cfg.Workspace == "" {
			return errors.NewGeneralError("no project, parent, or workspace specified", nil)
		}
		req.Workspace = cfg.Workspace
	}

	if cfg.DryRun {
		out := output.NewJSON(os.Stdout)
		return out.Print(map[string]any{"dry_run": true, "request": req})
	}

	client := newClient(cfg)
	task, err := client.CreateTask(context.Background(), req)
	if err != nil {
		return err
	}

	out := output.NewJSON(os.Stdout)
	return out.Print(task)
}

func runTaskUpdate(_ *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	if err := requireAuth(cfg); err != nil {
		return err
	}

	req := models.TaskUpdateRequest{}
	if taskUpdateName != "" {
		req.Name = &taskUpdateName
	}
	if taskUpdateNotes != "" {
		req.Notes = &taskUpdateNotes
	}
	if taskUpdateAssignee != "" {
		req.Assignee = &taskUpdateAssignee
	}
	if taskUpdateDueOn != "" {
		req.DueOn = &taskUpdateDueOn
	}

	if cfg.DryRun {
		out := output.NewJSON(os.Stdout)
		return out.Print(map[string]any{"dry_run": true, "gid": args[0], "request": req})
	}

	client := newClient(cfg)
	task, err := client.UpdateTask(context.Background(), args[0], req)
	if err != nil {
		return err
	}

	out := output.NewJSON(os.Stdout)
	return out.Print(task)
}

func runTaskComplete(_ *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	if err := requireAuth(cfg); err != nil {
		return err
	}

	completed := true
	req := models.TaskUpdateRequest{Completed: &completed}

	if cfg.DryRun {
		out := output.NewJSON(os.Stdout)
		return out.Print(map[string]any{"dry_run": true, "gid": args[0], "action": "complete"})
	}

	client := newClient(cfg)
	task, err := client.UpdateTask(context.Background(), args[0], req)
	if err != nil {
		return err
	}

	out := output.NewJSON(os.Stdout)
	return out.Print(task)
}

func runTaskDelete(_ *cobra.Command, args []string) error {
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
	if err := client.DeleteTask(context.Background(), args[0]); err != nil {
		return err
	}

	out := output.NewJSON(os.Stdout)
	return out.Print(map[string]any{"deleted": true, "gid": args[0]})
}

func runTaskAssign(_ *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	if err := requireAuth(cfg); err != nil {
		return err
	}

	assignee := args[1]
	req := models.TaskUpdateRequest{Assignee: &assignee}

	if cfg.DryRun {
		out := output.NewJSON(os.Stdout)
		return out.Print(map[string]any{"dry_run": true, "gid": args[0], "assignee": assignee})
	}

	client := newClient(cfg)
	task, err := client.UpdateTask(context.Background(), args[0], req)
	if err != nil {
		return err
	}

	out := output.NewJSON(os.Stdout)
	return out.Print(task)
}
