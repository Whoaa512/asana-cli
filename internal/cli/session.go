package cli

import (
	"context"
	"os"

	"github.com/spf13/cobra"

	"github.com/whoaa512/asana-cli/internal/config"
	"github.com/whoaa512/asana-cli/internal/errors"
	"github.com/whoaa512/asana-cli/internal/output"
	"github.com/whoaa512/asana-cli/internal/session"
)

var sessionCmd = &cobra.Command{
	Use:   "session",
	Short: "Manage work sessions",
	Long:  "Start, end, and log progress on work sessions that sync to Asana.",
}

var sessionStartCmd = &cobra.Command{
	Use:   "start [<task-gid>]",
	Short: "Start a work session",
	Long:  "Begin a work session on a task. Captures git branch and records start time.",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runSessionStart,
}

var sessionEndCmd = &cobra.Command{
	Use:   "end",
	Short: "End current work session",
	Long:  "End the current session and post a summary to Asana.",
	RunE:  runSessionEnd,
}

var sessionStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current session info",
	RunE:  runSessionStatus,
}

var sessionLogCmd = &cobra.Command{
	Use:   "log <message>",
	Short: "Add note to current session",
	Long:  "Add a progress note to the current session. Posted when session ends.",
	Args:  cobra.ExactArgs(1),
	RunE:  runSessionLog,
}

var (
	sessionStartForce bool
	sessionEndSummary string
	sessionEndDiscard bool
	sessionLogType    string
)

func init() {
	rootCmd.AddCommand(sessionCmd)
	sessionCmd.AddCommand(sessionStartCmd)
	sessionCmd.AddCommand(sessionEndCmd)
	sessionCmd.AddCommand(sessionStatusCmd)
	sessionCmd.AddCommand(sessionLogCmd)

	sessionStartCmd.Flags().BoolVar(&sessionStartForce, "force", false, "Force start, discarding existing session")

	sessionEndCmd.Flags().StringVar(&sessionEndSummary, "summary", "", "Additional summary text")
	sessionEndCmd.Flags().BoolVar(&sessionEndDiscard, "discard", false, "Discard session without posting to Asana")

	sessionLogCmd.Flags().StringVar(&sessionLogType, "type", "progress", "Log type: progress, decision, blocker")
}

func getSessionDir() (string, error) {
	root := session.GetRepoRoot()
	if root != "" {
		return root, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	fallback := home + "/.config/asana-cli"
	if err := os.MkdirAll(fallback, 0755); err != nil {
		return "", err
	}
	return fallback, nil
}

func runSessionStart(_ *cobra.Command, args []string) error {
	dir, err := getSessionDir()
	if err != nil {
		return errors.NewGeneralError("failed to determine session directory", err)
	}

	existing, err := session.Load(dir)
	if err != nil {
		return errors.NewGeneralError("failed to load existing session", err)
	}

	if existing != nil && !sessionStartForce {
		return errors.NewGeneralError("session already exists, use --force to override", nil)
	}

	localCtx, err := config.LoadLocalContext()
	if err != nil {
		return errors.NewGeneralError("failed to load context", err)
	}

	var taskGID string
	if len(args) > 0 {
		taskGID = args[0]
	} else {
		taskGID = localCtx.Task
	}

	if taskGID == "" {
		return errors.NewInvalidArgsError("task-gid required (provide as argument or set via 'ctx task')")
	}

	var opts []session.SessionOption
	if branch := session.GetCurrentBranch(); branch != "" {
		opts = append(opts, session.WithBranch(branch))
	}
	if repo := session.GetRepoName(); repo != "" {
		opts = append(opts, session.WithRepo(repo))
	}
	if localCtx.Project != "" {
		opts = append(opts, session.WithProject(localCtx.Project))
	}

	sess := session.New(taskGID, opts...)
	if err := sess.Save(dir); err != nil {
		return errors.NewGeneralError("failed to save session", err)
	}

	if err := updateContextTask(taskGID); err != nil {
		return errors.NewGeneralError("failed to update context", err)
	}

	out := output.NewJSON(os.Stdout)
	return out.Print(map[string]any{
		"task_gid":     sess.TaskGID,
		"started_at":   sess.StartedAt,
		"git_branch":   sess.StartBranch,
		"repo":         sess.Repo,
		"session_path": sess.Path(),
	})
}

func updateContextTask(taskGID string) error {
	localCtx, err := config.LoadLocalContext()
	if err != nil {
		return err
	}
	localCtx.Task = taskGID

	dir, err := config.FindContextFileDir()
	if err != nil {
		return err
	}
	return localCtx.Save(dir)
}

func runSessionEnd(_ *cobra.Command, _ []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	if err := requireAuth(cfg); err != nil {
		return err
	}

	dir, err := getSessionDir()
	if err != nil {
		return errors.NewGeneralError("failed to determine session directory", err)
	}

	sess, err := session.Load(dir)
	if err != nil {
		return errors.NewGeneralError("failed to load session", err)
	}
	if sess == nil {
		return errors.NewGeneralError("no active session", nil)
	}

	if sessionEndDiscard {
		if cfg.DryRun {
			out := output.NewJSON(os.Stdout)
			return out.Print(map[string]any{
				"dry_run":      true,
				"action":       "discard",
				"task_gid":     sess.TaskGID,
				"session_path": sess.Path(),
			})
		}
		if err := session.Delete(dir); err != nil {
			return errors.NewGeneralError("failed to delete session", err)
		}
		out := output.NewJSON(os.Stdout)
		return out.Print(map[string]any{
			"discarded":    true,
			"task_gid":     sess.TaskGID,
			"duration":     sess.FormatDuration(),
			"session_path": sess.Path(),
		})
	}

	if !sess.HasLogs() && sessionEndSummary == "" {
		if cfg.DryRun {
			out := output.NewJSON(os.Stdout)
			return out.Print(map[string]any{
				"dry_run":      true,
				"action":       "end_no_post",
				"task_gid":     sess.TaskGID,
				"reason":       "no logs or summary to post",
				"session_path": sess.Path(),
			})
		}
		if err := session.Delete(dir); err != nil {
			return errors.NewGeneralError("failed to delete session", err)
		}
		out := output.NewJSON(os.Stdout)
		return out.Print(map[string]any{
			"ended":        true,
			"task_gid":     sess.TaskGID,
			"duration":     sess.FormatDuration(),
			"posted":       false,
			"reason":       "no logs or summary to post",
			"session_path": sess.Path(),
		})
	}

	endBranch := session.GetCurrentBranch()
	summary := sess.FormatSummary(endBranch, sessionEndSummary)

	if cfg.DryRun {
		out := output.NewJSON(os.Stdout)
		return out.Print(map[string]any{
			"dry_run":      true,
			"action":       "end_and_post",
			"task_gid":     sess.TaskGID,
			"duration":     sess.FormatDuration(),
			"summary":      summary,
			"session_path": sess.Path(),
		})
	}

	client := newClient(cfg)
	story, err := client.AddComment(context.Background(), sess.TaskGID, summary)
	if err != nil {
		return errors.NewGeneralError("failed to post summary to Asana (session preserved, use --discard to clear)", err)
	}

	if err := session.Delete(dir); err != nil {
		return errors.NewGeneralError("failed to delete session", err)
	}

	result := map[string]any{
		"ended":        true,
		"task_gid":     sess.TaskGID,
		"duration":     sess.FormatDuration(),
		"posted":       true,
		"story_gid":    story.GID,
		"session_path": sess.Path(),
	}

	out := output.NewJSON(os.Stdout)
	return out.Print(result)
}

func runSessionStatus(_ *cobra.Command, _ []string) error {
	dir, err := getSessionDir()
	if err != nil {
		return errors.NewGeneralError("failed to determine session directory", err)
	}

	sess, err := session.Load(dir)
	if err != nil {
		return errors.NewGeneralError("failed to load session", err)
	}
	if sess == nil {
		out := output.NewJSON(os.Stdout)
		return out.Print(map[string]any{"active": false})
	}

	result := map[string]any{
		"active":       true,
		"task_gid":     sess.TaskGID,
		"started_at":   sess.StartedAt,
		"elapsed":      sess.FormatDuration(),
		"git_branch":   sess.StartBranch,
		"repo":         sess.Repo,
		"log_count":    len(sess.Logs),
		"logs":         sess.Logs,
		"stale":        sess.IsStale(),
		"session_path": sess.Path(),
	}

	out := output.NewJSON(os.Stdout)
	return out.Print(result)
}

func runSessionLog(_ *cobra.Command, args []string) error {
	dir, err := getSessionDir()
	if err != nil {
		return errors.NewGeneralError("failed to determine session directory", err)
	}

	sess, err := session.Load(dir)
	if err != nil {
		return errors.NewGeneralError("failed to load session", err)
	}
	if sess == nil {
		return errors.NewGeneralError("no active session, start one with 'session start'", nil)
	}

	validTypes := map[string]bool{"progress": true, "decision": true, "blocker": true}
	if !validTypes[sessionLogType] {
		return errors.NewInvalidArgsError("invalid log type, must be: progress, decision, or blocker")
	}

	sess.AddLog(sessionLogType, args[0])

	if err := sess.Save(dir); err != nil {
		return errors.NewGeneralError("failed to save session", err)
	}

	out := output.NewJSON(os.Stdout)
	return out.Print(map[string]any{
		"logged":       true,
		"type":         sessionLogType,
		"message":      args[0],
		"log_count":    len(sess.Logs),
		"session_path": sess.Path(),
	})
}
