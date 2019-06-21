package jokes

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"chucknorris/models"
	"chucknorris/pkg/api"
	"chucknorris/pkg/errors"
	e "chucknorris/pkg/errors"
)

const (
	jokesRandomOpt = "jokes/random/"
	jokesEndpoint  = "{{.Host}}/{{.Version}}"
)

// Service holds the rest client and configuration
type Service struct {
	api.Client
	Config api.Config
}

// NewService creates a new jokes service by passing the config
func NewService(rstCfg api.Config) *Service {
	svc := Service{
		Config: rstCfg,
		Client: rstCfg.New(),
	}
	return &svc
}

// JokesBuilder creates a GET requet that will call the jokes PI and return a Jokes response
func (svc *Service) JokesBuilder(ctx context.Context, req JokeReq) (*models.JokesResponse, error) {
	log.Println("performing JokesBuilder req")

	buf := &bytes.Buffer{}
	json.NewEncoder(buf).Encode(req)

	urlParams := api.URLParams{
		Host:    svc.Config.BaseURL,
		Version: jokesRandomOpt,
		Query: map[string][]string{
			"firstName": []string{req.FirstName},
			"lastName":  []string{req.LastName},
			"limitTo":   []string{req.LimitTo[0]},
		},
	}

	url, err := api.URLBuilder(svc.Config.Name, jokesEndpoint, urlParams)
	if err != nil {
		return nil, fmt.Errorf("not able to create API endpoint: %s", err)
	}

	var res models.JokesResponse

	// Perform request and return response or error
	if err := svc.Request(ctx, api.Request{
		Method: http.MethodGet,
		Body:   buf,
		URL:    url,
		Query:  urlParams.Query,
	}, &res); err != nil {
		// Check error type
		errorType := e.GetType(err)
		switch errorType {
		case e.BadRequest:
			// Return error message but don't kill app
			return nil, err
		case e.NotFound:
			return nil, err
		default:
			// Kill app because error is application level
			return nil, errors.Wrap(err, "joke request failed")
		}
	}
	return &res, nil
}
