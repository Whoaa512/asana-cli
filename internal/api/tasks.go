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
	OptFields []string
}

func (c *HTTPClient) ListTasks(ctx context.Context, opts TaskListOptions) (*models.ListResponse[models.Task], error) {
	var path string
	var isSearchEndpoint bool

	if opts.Tag != "" {
		path = fmt.Sprintf("/tags/%s/tasks", opts.Tag)
	} else if opts.Project != "" {
		path = fmt.Sprintf("/projects/%s/tasks", opts.Project)
	} else if opts.Workspace != "" {
		path = fmt.Sprintf("/workspaces/%s/tasks/search", opts.Workspace)
		isSearchEndpoint = true
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
		if isSearchEndpoint {
			path += fmt.Sprintf("%sassignee.any=%s", sep, opts.Assignee)
		} else {
			path += fmt.Sprintf("%sassignee=%s", sep, opts.Assignee)
		}
		sep = "&"
	}
	if opts.Completed != nil {
		path += fmt.Sprintf("%scompleted=%t", sep, *opts.Completed)
		sep = "&"
	}
	if len(opts.OptFields) > 0 {
		fields := ""
		for i, field := range opts.OptFields {
			if i > 0 {
				fields += ","
			}
			fields += field
		}
		path += fmt.Sprintf("%sopt_fields=%s", sep, fields)
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

func (c *HTTPClient) AddFollowers(ctx context.Context, taskGID string, followers []string) (*models.Task, error) {
	payload := struct {
		Data struct {
			Followers []string `json:"followers"`
		} `json:"data"`
	}{Data: struct {
		Followers []string `json:"followers"`
	}{Followers: followers}}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to encode request: %w", err)
	}

	var response struct {
		Data models.Task `json:"data"`
	}

	if err := c.post(ctx, fmt.Sprintf("/tasks/%s/addFollowers", taskGID), bytes.NewReader(body), &response); err != nil {
		return nil, err
	}

	return &response.Data, nil
}

func (c *HTTPClient) RemoveFollower(ctx context.Context, taskGID string, followerGID string) (*models.Task, error) {
	payload := struct {
		Data struct {
			Follower string `json:"follower"`
		} `json:"data"`
	}{Data: struct {
		Follower string `json:"follower"`
	}{Follower: followerGID}}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to encode request: %w", err)
	}

	var response struct {
		Data models.Task `json:"data"`
	}

	if err := c.post(ctx, fmt.Sprintf("/tasks/%s/removeFollower", taskGID), bytes.NewReader(body), &response); err != nil {
		return nil, err
	}

	return &response.Data, nil
}

func (c *HTTPClient) AddTag(ctx context.Context, taskGID string, tagGID string) (*models.Task, error) {
	payload := struct {
		Data struct {
			Tag string `json:"tag"`
		} `json:"data"`
	}{Data: struct {
		Tag string `json:"tag"`
	}{Tag: tagGID}}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to encode request: %w", err)
	}

	var response struct {
		Data models.Task `json:"data"`
	}

	if err := c.post(ctx, fmt.Sprintf("/tasks/%s/addTag", taskGID), bytes.NewReader(body), &response); err != nil {
		return nil, err
	}

	return &response.Data, nil
}

func (c *HTTPClient) RemoveTag(ctx context.Context, taskGID string, tagGID string) (*models.Task, error) {
	payload := struct {
		Data struct {
			Tag string `json:"tag"`
		} `json:"data"`
	}{Data: struct {
		Tag string `json:"tag"`
	}{Tag: tagGID}}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to encode request: %w", err)
	}

	var response struct {
		Data models.Task `json:"data"`
	}

	if err := c.post(ctx, fmt.Sprintf("/tasks/%s/removeTag", taskGID), bytes.NewReader(body), &response); err != nil {
		return nil, err
	}

	return &response.Data, nil
}

func (c *HTTPClient) AddToProject(ctx context.Context, taskGID string, projectGID string) (*models.Task, error) {
	payload := struct {
		Data struct {
			Project string `json:"project"`
		} `json:"data"`
	}{Data: struct {
		Project string `json:"project"`
	}{Project: projectGID}}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to encode request: %w", err)
	}

	var response struct {
		Data models.Task `json:"data"`
	}

	if err := c.post(ctx, fmt.Sprintf("/tasks/%s/addProject", taskGID), bytes.NewReader(body), &response); err != nil {
		return nil, err
	}

	return &response.Data, nil
}

func (c *HTTPClient) RemoveFromProject(ctx context.Context, taskGID string, projectGID string) (*models.Task, error) {
	payload := struct {
		Data struct {
			Project string `json:"project"`
		} `json:"data"`
	}{Data: struct {
		Project string `json:"project"`
	}{Project: projectGID}}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to encode request: %w", err)
	}

	var response struct {
		Data models.Task `json:"data"`
	}

	if err := c.post(ctx, fmt.Sprintf("/tasks/%s/removeProject", taskGID), bytes.NewReader(body), &response); err != nil {
		return nil, err
	}

	return &response.Data, nil
}

type TaskDuplicateRequest struct {
	Name    string   `json:"name,omitempty"`
	Include []string `json:"include,omitempty"`
}

func (c *HTTPClient) DuplicateTask(ctx context.Context, taskGID string, req TaskDuplicateRequest) (*models.Task, error) {
	payload := struct {
		Data TaskDuplicateRequest `json:"data"`
	}{Data: req}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to encode request: %w", err)
	}

	var response struct {
		Data models.Task `json:"data"`
	}

	if err := c.post(ctx, fmt.Sprintf("/tasks/%s/duplicate", taskGID), bytes.NewReader(body), &response); err != nil {
		return nil, err
	}

	return &response.Data, nil
}

func (c *HTTPClient) SetParent(ctx context.Context, taskGID string, parentGID *string) (*models.Task, error) {
	payload := struct {
		Data struct {
			Parent *string `json:"parent"`
		} `json:"data"`
	}{Data: struct {
		Parent *string `json:"parent"`
	}{Parent: parentGID}}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to encode request: %w", err)
	}

	var response struct {
		Data models.Task `json:"data"`
	}

	if err := c.post(ctx, fmt.Sprintf("/tasks/%s/setParent", taskGID), bytes.NewReader(body), &response); err != nil {
		return nil, err
	}

	return &response.Data, nil
}

func (c *HTTPClient) ListTaskProjects(ctx context.Context, taskGID string) ([]models.AsanaResource, error) {
	var response struct {
		Data []models.AsanaResource `json:"data"`
	}

	if err := c.get(ctx, fmt.Sprintf("/tasks/%s/projects", taskGID), &response); err != nil {
		return nil, err
	}

	return response.Data, nil
}
