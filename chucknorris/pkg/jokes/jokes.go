package jokes

import (
	"context"

	"chucknorris/models"
	"chucknorris/pkg/api"
)

// JokeReq contains the data needed to call Jokes API
type JokeReq struct {
	FirstName string   `json:"firstName"`
	LastName  string   `json:"lastName"`
	LimitTo   []string `json:"limitTo"`
}

// JokeService ...
type JokeService interface {
	api.Client
	JokesBuilder(ctx context.Context, req JokeReq) (*models.JokesResponse, error)
}
