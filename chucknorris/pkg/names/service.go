package names

import (
	"context"
	"log"
	"net/http"

	"chucknorris/models"
	"chucknorris/pkg/api"
	e "chucknorris/pkg/errors"

	"github.com/pkg/errors"
)

const (
	nameURL = "http://uinames.com/api/"
)

// Service holds the rest client and configuration
type Service struct {
	api.Client
	Config api.Config
}

// NewService creates a new names service by passing the config
func NewService(rstCfg api.Config) *Service {
	svc := Service{
		Config: rstCfg,
		Client: rstCfg.New(),
	}
	return &svc
}

// Names creates a GET request that will call the names API and returns a Name response
func (svc *Service) Names(ctx context.Context) (*models.NamesResponse, error) {
	log.Println("performing Names req")

	var res models.NamesResponse

	// Perform request and return response or error
	if err := svc.Request(ctx, api.Request{
		Method: http.MethodGet,
		URL:    nameURL,
	}, &res); err != nil {
		// check error type
		errorType := e.GetType(err)
		switch errorType {
		case e.BadRequest:
			// Return error message but don't kill app
			return nil, err
		case e.NotFound:
			return nil, err
		default:
			// Kill app because error is application level
			return nil, errors.Wrap(err, "name request failed")
		}
	}
	return &res, nil
}
