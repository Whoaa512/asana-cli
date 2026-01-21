package api

import (
	"context"
	"fmt"
	"net/url"
	"strings"

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
	OptFields []string
}

func (c *HTTPClient) SearchTasks(ctx context.Context, opts SearchTasksOptions) (*models.ListResponse[models.Task], error) {
	if opts.Workspace == "" {
		return nil, fmt.Errorf("workspace is required for search")
	}
	if opts.Text == "" {
		return nil, fmt.Errorf("text query is required for search")
	}

	path := fmt.Sprintf("/workspaces/%s/tasks/search", opts.Workspace)

	params := url.Values{}
	params.Set("text", opts.Text)
	if len(opts.OptFields) > 0 {
		params.Set("opt_fields", strings.Join(opts.OptFields, ","))
	}
	if opts.Limit > 0 {
		params.Set("limit", fmt.Sprintf("%d", opts.Limit))
	}
	if opts.Offset != "" {
		params.Set("offset", opts.Offset)
	}
	if opts.Project != "" {
		params.Set("projects.any", opts.Project)
	}
	if opts.Assignee != "" {
		params.Set("assignee.any", opts.Assignee)
	}
	if opts.Completed != nil {
		params.Set("completed", fmt.Sprintf("%t", *opts.Completed))
	}
	path += "?" + params.Encode()

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
