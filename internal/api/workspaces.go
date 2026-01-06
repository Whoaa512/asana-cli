package api

import (
	"context"
	"fmt"

	"github.com/whoaa512/asana-cli/internal/models"
)

func (c *HTTPClient) ListWorkspaces(ctx context.Context, limit int) (*models.ListResponse[models.Workspace], error) {
	path := "/workspaces"
	if limit > 0 {
		path = fmt.Sprintf("%s?limit=%d", path, limit)
	}

	var response struct {
		Data     []models.Workspace `json:"data"`
		NextPage *models.PageInfo   `json:"next_page,omitempty"`
	}

	if err := c.get(ctx, path, &response); err != nil {
		return nil, err
	}

	return &models.ListResponse[models.Workspace]{
		Data:     response.Data,
		NextPage: response.NextPage,
	}, nil
}

func (c *HTTPClient) GetWorkspace(ctx context.Context, gid string) (*models.Workspace, error) {
	var response struct {
		Data models.Workspace `json:"data"`
	}

	if err := c.get(ctx, "/workspaces/"+gid, &response); err != nil {
		return nil, err
	}

	return &response.Data, nil
}
