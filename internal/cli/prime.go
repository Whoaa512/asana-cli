package cli

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/whoaa512/asana-cli/internal/api"
	"github.com/whoaa512/asana-cli/internal/errors"
	"github.com/whoaa512/asana-cli/internal/models"
	"github.com/whoaa512/asana-cli/internal/session"
)

var primeCmd = &cobra.Command{
	Use:   "prime",
	Short: "Output markdown context dump for AI agents",
	Long: `Output a markdown context dump of your current Asana state.

Includes:
- Active session (if exists)
- Ready tasks (unblocked)
- Blocked tasks
- Task subtasks (1 level)
- Recent activity (last 24h)

Output is raw markdown to stdout for easy piping to AI agents.`,
	RunE: runPrime,
}

var (
	primeProject          string
	primeLimit            int
	primeIncludeCompleted bool
)

func init() {
	rootCmd.AddCommand(primeCmd)
	primeCmd.Flags().StringVar(&primeProject, "project", "", "Override project GID (default from context)")
	primeCmd.Flags().IntVar(&primeLimit, "limit", 20, "Max tasks to show per section")
	primeCmd.Flags().BoolVar(&primeIncludeCompleted, "include-completed", false, "Show recently completed tasks")
}

func runPrime(_ *cobra.Command, _ []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	if err := requireAuth(cfg); err != nil {
		return err
	}

	project := primeProject
	if project == "" && cfg.Project != "" {
		project = cfg.Project
	}
	if project == "" {
		return errors.NewGeneralError("no project specified via --project or context", nil)
	}

	if cfg.DryRun {
		fmt.Fprintln(os.Stderr, "# Dry run - would fetch:")
		fmt.Fprintf(os.Stderr, "- Project: %s\n", project)
		fmt.Fprintf(os.Stderr, "- Limit: %d\n", primeLimit)
		fmt.Fprintf(os.Stderr, "- Include completed: %v\n", primeIncludeCompleted)
		return nil
	}

	client := newClient(cfg)
	ctx := context.Background()

	var output strings.Builder
	output.WriteString("# Asana Context\n\n")

	if err := writeActiveSession(&output, client, ctx); err != nil {
		return err
	}

	readyTasks, blockedTasks, err := fetchAndCategorize(client, project, primeLimit)
	if err != nil {
		return err
	}

	if err := writeReadyTasks(&output, client, ctx, readyTasks); err != nil {
		return err
	}

	if err := writeBlockedTasks(&output, client, ctx, blockedTasks); err != nil {
		return err
	}

	if primeIncludeCompleted {
		completedTasks, err := fetchCompletedTasks(client, project, primeLimit)
		if err != nil {
			return err
		}
		writeCompletedTasks(&output, completedTasks)
	}

	fmt.Print(output.String())
	return nil
}

func writeActiveSession(output *strings.Builder, client api.Client, ctx context.Context) error {
	dir, err := getSessionDir()
	if err != nil {
		return nil
	}

	sess, err := session.Load(dir)
	if err != nil || sess == nil {
		return nil
	}

	task, err := client.GetTask(ctx, sess.TaskGID)
	if err != nil {
		return nil
	}

	output.WriteString("## Active Session\n")
	output.WriteString(fmt.Sprintf("Task: %s (%s)\n", task.Name, sess.TaskGID))
	output.WriteString(fmt.Sprintf("Started: %s ago\n", sess.FormatDuration()))

	if sess.StartBranch != "" {
		currentBranch := session.GetCurrentBranch()
		if currentBranch != "" && currentBranch != sess.StartBranch {
			output.WriteString(fmt.Sprintf("Branch: %s → %s\n", sess.StartBranch, currentBranch))
		} else {
			output.WriteString(fmt.Sprintf("Branch: %s\n", sess.StartBranch))
		}
	}

	if len(sess.Logs) > 0 {
		output.WriteString("Progress logs:\n")
		for _, log := range sess.Logs {
			output.WriteString(fmt.Sprintf("- [%s] %s\n", log.Timestamp.Local().Format("15:04"), log.Text))
		}
	}

	stories, err := client.ListStories(ctx, sess.TaskGID, 100, "")
	if err == nil && len(stories.Data) > 0 {
		output.WriteString("\nRecent comments (last 24h):\n")
		cutoff := time.Now().Add(-24 * time.Hour)
		hasRecent := false
		for _, story := range stories.Data {
			if story.Type == "comment" && story.Text != "" {
				createdAt, err := time.Parse(time.RFC3339, story.CreatedAt)
				if err == nil && createdAt.After(cutoff) {
					hasRecent = true
					by := "Unknown"
					if story.CreatedBy != nil {
						by = story.CreatedBy.Name
					}
					output.WriteString(fmt.Sprintf("- [%s] %s: %s\n",
						createdAt.Local().Format("15:04"),
						by,
						truncate(story.Text, 100)))
				}
			}
		}
		if !hasRecent {
			output.WriteString("- (none)\n")
		}
	}

	output.WriteString("\n")
	return nil
}

func writeReadyTasks(output *strings.Builder, client api.Client, ctx context.Context, tasks []models.Task) error {
	if len(tasks) == 0 {
		return nil
	}

	output.WriteString("## Ready Tasks (unblocked)\n")
	for _, task := range tasks {
		dueStr := ""
		if task.DueOn != "" {
			dueStr = fmt.Sprintf(" - due %s", task.DueOn)
		}
		output.WriteString(fmt.Sprintf("- [ ] %s (%s)%s\n", task.Name, task.GID, dueStr))

		subtasks, err := client.ListSubtasks(ctx, task.GID, 50, "")
		if err == nil && len(subtasks.Data) > 0 {
			for _, subtask := range subtasks.Data {
				checkbox := "[ ]"
				if subtask.Completed {
					checkbox = "[x]"
				}
				output.WriteString(fmt.Sprintf("  - %s %s\n", checkbox, subtask.Name))
			}
		}
	}
	output.WriteString("\n")
	return nil
}

func writeBlockedTasks(output *strings.Builder, client api.Client, ctx context.Context, tasks []models.Task) error {
	if len(tasks) == 0 {
		return nil
	}

	output.WriteString("## Blocked Tasks\n")
	for _, task := range tasks {
		deps, err := client.ListDependencies(ctx, task.GID)
		if err != nil {
			continue
		}

		var blockers []string
		for _, dep := range deps {
			if !dep.Completed {
				blockers = append(blockers, dep.Name)
			}
		}

		blockStr := ""
		if len(blockers) > 0 {
			blockStr = fmt.Sprintf(" - blocked by: %s", strings.Join(blockers, ", "))
		}
		output.WriteString(fmt.Sprintf("- [ ] %s (%s)%s\n", task.Name, task.GID, blockStr))

		subtasks, err := client.ListSubtasks(ctx, task.GID, 50, "")
		if err == nil && len(subtasks.Data) > 0 {
			for _, subtask := range subtasks.Data {
				checkbox := "[ ]"
				if subtask.Completed {
					checkbox = "[x]"
				}
				output.WriteString(fmt.Sprintf("  - %s %s\n", checkbox, subtask.Name))
			}
		}
	}
	output.WriteString("\n")
	return nil
}

func fetchAndCategorize(client api.Client, project string, limit int) ([]models.Task, []models.Task, error) {
	incompleteTasks, err := fetchIncompleteTasks(client, project, "", limit)
	if err != nil {
		return nil, nil, err
	}

	var ready []models.Task
	var blocked []models.Task

	for _, task := range incompleteTasks {
		hasBlockers, err := hasIncompleteDependencies(client, task.GID)
		if err != nil {
			return nil, nil, err
		}

		if hasBlockers {
			blocked = append(blocked, task)
		} else {
			ready = append(ready, task)
		}
	}

	return ready, blocked, nil
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

func fetchCompletedTasks(client api.Client, project string, limit int) ([]models.Task, error) {
	completed := true
	opts := api.TaskListOptions{
		Project:   project,
		Completed: &completed,
		Limit:     limit,
	}

	result, err := client.ListTasks(context.Background(), opts)
	if err != nil {
		return nil, err
	}

	return result.Data, nil
}

func writeCompletedTasks(output *strings.Builder, tasks []models.Task) {
	if len(tasks) == 0 {
		return
	}

	output.WriteString("## Recently Completed\n")
	for _, task := range tasks {
		completedOn := ""
		if task.CompletedAt != "" {
			completedOn = fmt.Sprintf(" - completed %s", task.CompletedAt[:10])
		}
		output.WriteString(fmt.Sprintf("- [x] %s (%s)%s\n", task.Name, task.GID, completedOn))
	}
	output.WriteString("\n")
}
