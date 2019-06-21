package processor

import (
	"context"

	"chucknorris/pkg/api"
	"chucknorris/pkg/jokes"
	"chucknorris/pkg/names"
)

// JokeService ...
type JokeService struct {
	cfg     Configs
	clients APIClients
}

// APIClients contains the API services we will use
type APIClients struct {
	Jokes jokes.JokeService
	Names names.NameService
}

// Configs ...
type Configs struct {
	JokesCfg api.Config `yaml:"joke_service" json:"joke_service" mapstructure:"joke_service"`
	NamesCfg api.Config `yaml:"name_service" json:"name_service" mapstructure:"name_service"`
}

// New returns a new JokeService
func New(ctx context.Context, apis APIClients, cfg Configs) *JokeService {
	js := &JokeService{
		cfg:     cfg,
		clients: apis,
	}
	return js
}
