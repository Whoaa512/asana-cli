package api

import (
	"context"
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
