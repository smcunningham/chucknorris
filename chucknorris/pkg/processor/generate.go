package processor

import (
	"context"
	"fmt"
	"log"

	"chucknorris/models"
	"chucknorris/pkg/errors"
	e "chucknorris/pkg/errors"
	"chucknorris/pkg/jokes"
)

// GenerateName will use the client to call the Name API and return its response
func (js *JokeService) GenerateName(ctx context.Context) (*models.NamesResponse, error) {
	log.Println("generating name")
	name, err := js.clients.Names.Names(ctx)
	if err != nil {
		errorType := e.GetType(err)
		switch errorType {
		case e.BadRequest:
			// Return error message but don't kill app
			return nil, err
		case e.NotFound:
			return nil, err
		default:
			// Kill app because error is application level
			return nil, errors.Wrap(err, "name generate failed")
		}
	}
	fmt.Printf("\nName: %s %s\n", name.Name, name.Surname)

	return name, nil
}

// GenerateJoke will use the client to call the Jokes API and return a joke as a string from the response
func (js *JokeService) GenerateJoke(ctx context.Context, req jokes.JokeReq) (string, error) {
	log.Println("generating joke")
	res, err := js.clients.Jokes.JokesBuilder(ctx, req)
	if err != nil {
		errorType := e.GetType(err)
		switch errorType {
		case e.BadRequest:
			// Return error message but don't kill app
			return "", err
		case e.NotFound:
			return "", err
		default:
			// Kill app because error is application level
			return "", errors.Wrap(err, "joke generate failed")
		}
	}

	return res.Value.Joke, nil
}
