package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"

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

	params := url.Values{}
	if opts.Limit > 0 {
		params.Set("limit", fmt.Sprintf("%d", opts.Limit))
	}
	if opts.Offset != "" {
		params.Set("offset", opts.Offset)
	}
	if opts.Archived {
		params.Set("archived", fmt.Sprintf("%t", opts.Archived))
	}
	if len(params) > 0 {
		path += "?" + params.Encode()
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

	path := fmt.Sprintf("/workspaces/%s/projects", opts.Workspace)

	params := url.Values{}
	params.Set("opt_fields", "name,archived,workspace,color")
	if opts.Limit > 0 {
		params.Set("limit", fmt.Sprintf("%d", opts.Limit))
	}
	if opts.Offset != "" {
		params.Set("offset", opts.Offset)
	}
	path += "?" + params.Encode()

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
