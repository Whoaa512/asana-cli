package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/whoaa512/asana-cli/internal/models"
)

type TagListOptions struct {
	Workspace string
	Limit     int
	Offset    string
}

func (c *HTTPClient) ListTags(ctx context.Context, opts TagListOptions) (*models.ListResponse[models.Tag], error) {
	if opts.Workspace == "" {
		return nil, fmt.Errorf("workspace is required")
	}

	path := fmt.Sprintf("/workspaces/%s/tags", opts.Workspace)
	sep := "?"

	if opts.Limit > 0 {
		path += fmt.Sprintf("%slimit=%d", sep, opts.Limit)
		sep = "&"
	}
	if opts.Offset != "" {
		path += fmt.Sprintf("%soffset=%s", sep, opts.Offset)
	}

	var response struct {
		Data     []models.Tag     `json:"data"`
		NextPage *models.PageInfo `json:"next_page,omitempty"`
	}

	if err := c.get(ctx, path, &response); err != nil {
		return nil, err
	}

	return &models.ListResponse[models.Tag]{
		Data:     response.Data,
		NextPage: response.NextPage,
	}, nil
}

func (c *HTTPClient) GetTag(ctx context.Context, gid string) (*models.Tag, error) {
	var response struct {
		Data models.Tag `json:"data"`
	}

	if err := c.get(ctx, "/tags/"+gid, &response); err != nil {
		return nil, err
	}

	return &response.Data, nil
}

type TagCreateRequest struct {
	Name      string
	Workspace string
	Color     string
}

func (c *HTTPClient) CreateTag(ctx context.Context, req TagCreateRequest) (*models.Tag, error) {
	if req.Workspace == "" {
		return nil, fmt.Errorf("workspace is required")
	}
	if req.Name == "" {
		return nil, fmt.Errorf("name is required")
	}

	path := fmt.Sprintf("/workspaces/%s/tags", req.Workspace)

	payload := map[string]any{
		"data": map[string]any{
			"name": req.Name,
		},
	}
	if req.Color != "" {
		payload["data"].(map[string]any)["color"] = req.Color
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	var response struct {
		Data models.Tag `json:"data"`
	}

	if err := c.post(ctx, path, bytes.NewReader(body), &response); err != nil {
		return nil, err
	}

	return &response.Data, nil
}
