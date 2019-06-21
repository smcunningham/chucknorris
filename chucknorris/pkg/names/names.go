package names

import (
	"context"

	"chucknorris/models"
	"chucknorris/pkg/api"
)

// NameService is the interface for names through the name API
// Names retrieves the names
type NameService interface {
	api.Client
	Names(ctx context.Context) (*models.NamesResponse, error)
}
