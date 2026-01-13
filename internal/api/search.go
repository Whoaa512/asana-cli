package api

import (
	"context"
	"fmt"

	"github.com/whoaa512/asana-cli/internal/models"
)

type SearchTasksOptions struct {
	Workspace string
	Text      string
	Project   string
	Assignee  string
	Completed *bool
	Limit     int
	Offset    string
}

func (c *HTTPClient) SearchTasks(ctx context.Context, opts SearchTasksOptions) (*models.ListResponse[models.Task], error) {
	if opts.Workspace == "" {
		return nil, fmt.Errorf("workspace is required for search")
	}
	if opts.Text == "" {
		return nil, fmt.Errorf("text query is required for search")
	}

	path := fmt.Sprintf("/workspaces/%s/tasks/search", opts.Workspace)
	sep := "?"

	path += fmt.Sprintf("%stext=%s", sep, opts.Text)
	sep = "&"

	if opts.Limit > 0 {
		path += fmt.Sprintf("%slimit=%d", sep, opts.Limit)
		sep = "&"
	}
	if opts.Offset != "" {
		path += fmt.Sprintf("%soffset=%s", sep, opts.Offset)
		sep = "&"
	}
	if opts.Project != "" {
		path += fmt.Sprintf("%sprojects.any=%s", sep, opts.Project)
		sep = "&"
	}
	if opts.Assignee != "" {
		path += fmt.Sprintf("%sassignee.any=%s", sep, opts.Assignee)
		sep = "&"
	}
	if opts.Completed != nil {
		path += fmt.Sprintf("%scompleted=%t", sep, *opts.Completed)
	}

	var response struct {
		Data     []models.Task    `json:"data"`
		NextPage *models.PageInfo `json:"next_page,omitempty"`
	}

	if err := c.get(ctx, path, &response); err != nil {
		return nil, err
	}

	return &models.ListResponse[models.Task]{
		Data:     response.Data,
		NextPage: response.NextPage,
	}, nil
}
