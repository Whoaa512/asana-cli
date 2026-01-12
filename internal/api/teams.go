package api

import (
	"context"
	"fmt"

	"github.com/whoaa512/asana-cli/internal/models"
)

type TeamListOptions struct {
	Organization string
	Limit        int
	Offset       string
}

func (c *HTTPClient) ListTeams(ctx context.Context, opts TeamListOptions) (*models.ListResponse[models.Team], error) {
	if opts.Organization == "" {
		return nil, fmt.Errorf("organization is required")
	}

	path := fmt.Sprintf("/organizations/%s/teams?opt_fields=name", opts.Organization)

	if opts.Limit > 0 {
		path += fmt.Sprintf("&limit=%d", opts.Limit)
	}
	if opts.Offset != "" {
		path += fmt.Sprintf("&offset=%s", opts.Offset)
	}

	var response struct {
		Data     []models.Team    `json:"data"`
		NextPage *models.PageInfo `json:"next_page,omitempty"`
	}

	if err := c.get(ctx, path, &response); err != nil {
		return nil, err
	}

	return &models.ListResponse[models.Team]{
		Data:     response.Data,
		NextPage: response.NextPage,
	}, nil
}
