package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/whoaa512/asana-cli/internal/models"
)

func (c *HTTPClient) ListStories(ctx context.Context, taskGID string, limit int, offset string) (*models.ListResponse[models.Story], error) {
	path := fmt.Sprintf("/tasks/%s/stories", taskGID)

	sep := "?"
	if limit > 0 {
		path += fmt.Sprintf("%slimit=%d", sep, limit)
		sep = "&"
	}
	if offset != "" {
		path += fmt.Sprintf("%soffset=%s", sep, offset)
	}

	var response struct {
		Data     []models.Story   `json:"data"`
		NextPage *models.PageInfo `json:"next_page,omitempty"`
	}

	if err := c.get(ctx, path, &response); err != nil {
		return nil, err
	}

	return &models.ListResponse[models.Story]{
		Data:     response.Data,
		NextPage: response.NextPage,
	}, nil
}

func (c *HTTPClient) AddComment(ctx context.Context, taskGID string, text string) (*models.Story, error) {
	payload := struct {
		Data models.StoryCreateRequest `json:"data"`
	}{Data: models.StoryCreateRequest{Text: text}}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to encode request: %w", err)
	}

	var response struct {
		Data models.Story `json:"data"`
	}

	if err := c.post(ctx, fmt.Sprintf("/tasks/%s/stories", taskGID), bytes.NewReader(body), &response); err != nil {
		return nil, err
	}

	return &response.Data, nil
}
