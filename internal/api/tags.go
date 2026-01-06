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

func (c *HTTPClient) CreateTag(ctx context.Context, workspaceGID string, req models.TagCreateRequest) (*models.Tag, error) {
	payload := struct {
		Data struct {
			Name      string `json:"name"`
			Color     string `json:"color,omitempty"`
			Workspace string `json:"workspace"`
		} `json:"data"`
	}{}
	payload.Data.Name = req.Name
	payload.Data.Color = req.Color
	payload.Data.Workspace = workspaceGID

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to encode request: %w", err)
	}

	var response struct {
		Data models.Tag `json:"data"`
	}

	if err := c.post(ctx, "/tags", bytes.NewReader(body), &response); err != nil {
		return nil, err
	}

	return &response.Data, nil
}
