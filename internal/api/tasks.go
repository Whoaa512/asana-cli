package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/whoaa512/asana-cli/internal/models"
)

type TaskListOptions struct {
	Project   string
	Assignee  string
	Completed *bool
	Limit     int
	Offset    string
	Workspace string
	Tag       string
}

func (c *HTTPClient) ListTasks(ctx context.Context, opts TaskListOptions) (*models.ListResponse[models.Task], error) {
	var path string

	if opts.Tag != "" {
		path = fmt.Sprintf("/tags/%s/tasks", opts.Tag)
	} else if opts.Project != "" {
		path = fmt.Sprintf("/projects/%s/tasks", opts.Project)
	} else if opts.Workspace != "" {
		path = fmt.Sprintf("/workspaces/%s/tasks/search", opts.Workspace)
	} else {
		return nil, fmt.Errorf("either project, tag, or workspace is required")
	}

	sep := "?"
	if opts.Limit > 0 {
		path += fmt.Sprintf("%slimit=%d", sep, opts.Limit)
		sep = "&"
	}
	if opts.Offset != "" {
		path += fmt.Sprintf("%soffset=%s", sep, opts.Offset)
		sep = "&"
	}
	if opts.Assignee != "" {
		path += fmt.Sprintf("%sassignee=%s", sep, opts.Assignee)
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

func (c *HTTPClient) GetTask(ctx context.Context, gid string) (*models.Task, error) {
	var response struct {
		Data models.Task `json:"data"`
	}

	if err := c.get(ctx, "/tasks/"+gid, &response); err != nil {
		return nil, err
	}

	return &response.Data, nil
}

func (c *HTTPClient) CreateTask(ctx context.Context, req models.TaskCreateRequest) (*models.Task, error) {
	payload := struct {
		Data models.TaskCreateRequest `json:"data"`
	}{Data: req}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to encode request: %w", err)
	}

	var response struct {
		Data models.Task `json:"data"`
	}

	if err := c.post(ctx, "/tasks", bytes.NewReader(body), &response); err != nil {
		return nil, err
	}

	return &response.Data, nil
}

func (c *HTTPClient) UpdateTask(ctx context.Context, gid string, req models.TaskUpdateRequest) (*models.Task, error) {
	payload := struct {
		Data models.TaskUpdateRequest `json:"data"`
	}{Data: req}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to encode request: %w", err)
	}

	var response struct {
		Data models.Task `json:"data"`
	}

	if err := c.put(ctx, "/tasks/"+gid, bytes.NewReader(body), &response); err != nil {
		return nil, err
	}

	return &response.Data, nil
}

func (c *HTTPClient) DeleteTask(ctx context.Context, gid string) error {
	return c.delete(ctx, "/tasks/"+gid)
}
