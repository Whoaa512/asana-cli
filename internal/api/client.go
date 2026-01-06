package api

import (
	"context"

	"github.com/whoaa512/asana-cli/internal/models"
)

const BaseURL = "https://app.asana.com/api/1.0"

type Client interface {
	GetMe(ctx context.Context) (*models.User, error)
}
