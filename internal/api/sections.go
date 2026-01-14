package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/whoaa512/asana-cli/internal/models"
)

type SectionListOptions struct {
	Project string
	Limit   int
	Offset  string
}

func (c *HTTPClient) ListSections(ctx context.Context, opts SectionListOptions) (*models.ListResponse[models.Section], error) {
	if opts.Project == "" {
		return nil, fmt.Errorf("project is required")
	}

	path := fmt.Sprintf("/projects/%s/sections", opts.Project)
	sep := "?"

	if opts.Limit > 0 {
		path += fmt.Sprintf("%slimit=%d", sep, opts.Limit)
		sep = "&"
	}
	if opts.Offset != "" {
		path += fmt.Sprintf("%soffset=%s", sep, opts.Offset)
	}

	var response struct {
		Data     []models.Section `json:"data"`
		NextPage *models.PageInfo `json:"next_page,omitempty"`
	}

	if err := c.get(ctx, path, &response); err != nil {
		return nil, err
	}

	return &models.ListResponse[models.Section]{
		Data:     response.Data,
		NextPage: response.NextPage,
	}, nil
}

func (c *HTTPClient) GetSection(ctx context.Context, gid string) (*models.Section, error) {
	var response struct {
		Data models.Section `json:"data"`
	}

	if err := c.get(ctx, "/sections/"+gid, &response); err != nil {
		return nil, err
	}

	return &response.Data, nil
}

func (c *HTTPClient) CreateSection(ctx context.Context, projectGID string, req models.SectionCreateRequest) (*models.Section, error) {
	payload := struct {
		Data struct {
			Name string `json:"name"`
		} `json:"data"`
	}{}
	payload.Data.Name = req.Name

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to encode request: %w", err)
	}

	var response struct {
		Data models.Section `json:"data"`
	}

	path := fmt.Sprintf("/projects/%s/sections", projectGID)
	if err := c.post(ctx, path, bytes.NewReader(body), &response); err != nil {
		return nil, err
	}

	return &response.Data, nil
}

func (c *HTTPClient) AddTaskToSection(ctx context.Context, sectionGID string, taskGID string) error {
	payload := struct {
		Data struct {
			Task string `json:"task"`
		} `json:"data"`
	}{}
	payload.Data.Task = taskGID

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to encode request: %w", err)
	}

	path := fmt.Sprintf("/sections/%s/addTask", sectionGID)
	return c.post(ctx, path, bytes.NewReader(body), nil)
}

func (c *HTTPClient) UpdateSection(ctx context.Context, gid string, req models.SectionUpdateRequest) (*models.Section, error) {
	payload := struct {
		Data struct {
			Name string `json:"name"`
		} `json:"data"`
	}{}
	payload.Data.Name = req.Name

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to encode request: %w", err)
	}

	var response struct {
		Data models.Section `json:"data"`
	}

	if err := c.put(ctx, "/sections/"+gid, bytes.NewReader(body), &response); err != nil {
		return nil, err
	}

	return &response.Data, nil
}

func (c *HTTPClient) DeleteSection(ctx context.Context, gid string) error {
	return c.delete(ctx, "/sections/"+gid)
}

func (c *HTTPClient) InsertSection(ctx context.Context, projectGID string, req models.SectionInsertRequest) error {
	payload := struct {
		Data models.SectionInsertRequest `json:"data"`
	}{}
	payload.Data = req

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to encode request: %w", err)
	}

	path := fmt.Sprintf("/projects/%s/sections/insert", projectGID)
	return c.post(ctx, path, bytes.NewReader(body), nil)
}
