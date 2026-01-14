package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/whoaa512/asana-cli/internal/models"
)

type ProjectListOptions struct {
	Workspace string
	Archived  bool
	Limit     int
	Offset    string
}

type UserProjectListOptions struct {
	Workspace string
	Limit     int
	Offset    string
}

func (c *HTTPClient) ListProjects(ctx context.Context, opts ProjectListOptions) (*models.ListResponse[models.Project], error) {
	if opts.Workspace == "" {
		return nil, fmt.Errorf("workspace is required")
	}

	path := fmt.Sprintf("/workspaces/%s/projects", opts.Workspace)
	sep := "?"

	if opts.Limit > 0 {
		path += fmt.Sprintf("%slimit=%d", sep, opts.Limit)
		sep = "&"
	}
	if opts.Offset != "" {
		path += fmt.Sprintf("%soffset=%s", sep, opts.Offset)
		sep = "&"
	}
	if opts.Archived {
		path += fmt.Sprintf("%sarchived=%t", sep, opts.Archived)
	}

	var response struct {
		Data     []models.Project `json:"data"`
		NextPage *models.PageInfo `json:"next_page,omitempty"`
	}

	if err := c.get(ctx, path, &response); err != nil {
		return nil, err
	}

	return &models.ListResponse[models.Project]{
		Data:     response.Data,
		NextPage: response.NextPage,
	}, nil
}

func (c *HTTPClient) GetProject(ctx context.Context, gid string) (*models.Project, error) {
	var response struct {
		Data models.Project `json:"data"`
	}

	if err := c.get(ctx, "/projects/"+gid, &response); err != nil {
		return nil, err
	}

	return &response.Data, nil
}

func (c *HTTPClient) CreateProject(ctx context.Context, req models.ProjectCreateRequest) (*models.Project, error) {
	payload := struct {
		Data models.ProjectCreateRequest `json:"data"`
	}{Data: req}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to encode request: %w", err)
	}

	var response struct {
		Data models.Project `json:"data"`
	}

	if err := c.post(ctx, "/projects", bytes.NewReader(body), &response); err != nil {
		return nil, err
	}

	return &response.Data, nil
}

func (c *HTTPClient) ListUserProjects(ctx context.Context, opts UserProjectListOptions) (*models.ListResponse[models.Project], error) {
	if opts.Workspace == "" {
		return nil, fmt.Errorf("workspace is required")
	}

	path := fmt.Sprintf("/workspaces/%s/projects?opt_fields=name,archived,workspace,color", opts.Workspace)
	sep := "&"

	if opts.Limit > 0 {
		path += fmt.Sprintf("%slimit=%d", sep, opts.Limit)
	}
	if opts.Offset != "" {
		path += fmt.Sprintf("%soffset=%s", sep, opts.Offset)
	}

	var response struct {
		Data     []models.Project `json:"data"`
		NextPage *models.PageInfo `json:"next_page,omitempty"`
	}

	if err := c.get(ctx, path, &response); err != nil {
		return nil, err
	}

	return &models.ListResponse[models.Project]{
		Data:     response.Data,
		NextPage: response.NextPage,
	}, nil
}
