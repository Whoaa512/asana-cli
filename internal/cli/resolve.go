package cli

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/sahilm/fuzzy"

	"github.com/whoaa512/asana-cli/internal/api"
	"github.com/whoaa512/asana-cli/internal/config"
	"github.com/whoaa512/asana-cli/internal/errors"
	"github.com/whoaa512/asana-cli/internal/models"
)

var gidRegex = regexp.MustCompile(`^\d+$`)

type taskMatch struct {
	task  models.Task
	score int
}

func resolveTaskGID(ctx context.Context, cfg *config.Config, client api.Client, nameOrGID string, allowPick bool) (string, error) {
	if gidRegex.MatchString(nameOrGID) {
		return nameOrGID, nil
	}

	tasks, err := fetchRecentTasks(ctx, cfg, client)
	if err != nil {
		return "", err
	}

	matches := fuzzyMatchTasks(tasks, nameOrGID)

	if len(matches) == 0 {
		return "", errors.NewGeneralError(fmt.Sprintf("no tasks found matching '%s'", nameOrGID), nil)
	}

	if len(matches) == 1 {
		return matches[0].task.GID, nil
	}

	if !allowPick {
		return "", errors.NewGeneralError(fmt.Sprintf("multiple tasks match '%s', use --pick flag for interactive selection", nameOrGID), nil)
	}

	selected, err := pickTask(matches)
	if err != nil {
		return "", err
	}

	return selected.task.GID, nil
}

func fetchRecentTasks(ctx context.Context, cfg *config.Config, client api.Client) ([]models.Task, error) {
	opts := api.TaskListOptions{
		Assignee: "me",
		Limit:    200,
	}

	if cfg.Project != "" {
		opts.Project = cfg.Project
	} else if cfg.Workspace != "" {
		opts.Workspace = cfg.Workspace
	} else {
		return nil, errors.NewGeneralError("no project or workspace configured", nil)
	}

	result, err := client.ListTasks(ctx, opts)
	if err != nil {
		return nil, err
	}

	return result.Data, nil
}

func fuzzyMatchTasks(tasks []models.Task, query string) []taskMatch {
	queryLower := strings.ToLower(query)
	var matches []taskMatch

	taskNames := make([]string, len(tasks))
	for i, task := range tasks {
		taskNames[i] = task.Name
	}

	fuzzyResults := fuzzy.Find(query, taskNames)

	if len(fuzzyResults) > 0 {
		for _, result := range fuzzyResults {
			matches = append(matches, taskMatch{
				task:  tasks[result.Index],
				score: result.Score,
			})
		}
		return matches
	}

	for _, task := range tasks {
		if strings.Contains(strings.ToLower(task.Name), queryLower) {
			matches = append(matches, taskMatch{
				task:  task,
				score: 0,
			})
		}
	}

	return matches
}
