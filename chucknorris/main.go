package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"

	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"chucknorris/cmd/config"
	"chucknorris/pkg/api"
	e "chucknorris/pkg/errors"
	"chucknorris/pkg/jokes"
	"chucknorris/pkg/names"
	"chucknorris/pkg/processor"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
	}()

	// Read in configuration from .yaml file
	cfg, err := initConfig()
	if err != nil {
		log.Fatalf("error initiaizing configuration: %s", err)
	}

	// Name service configuration
	nameCfg := api.Config{
		Name:    cfg.Clients.Names.Name,
		BaseURL: cfg.Clients.Names.BaseURL,
		Timeout: cfg.Clients.Names.Timeout,
	}

	// Jokes service configuration
	jokesCfg := api.Config{
		Name:    cfg.Clients.Jokes.Name,
		BaseURL: cfg.Clients.Jokes.BaseURL,
		Timeout: cfg.Clients.Jokes.Timeout,
	}

	// Processor configuration
	svcCfg := processor.Configs{
		JokesCfg: jokesCfg,
		NamesCfg: nameCfg,
	}

	// Create new services using configurations
	nameSvc := names.NewService(nameCfg)
	jokesSvc := jokes.NewService(jokesCfg)

	apiClients := processor.APIClients{
		Jokes: jokesSvc,
		Names: nameSvc,
	}

	// Create processor using context, api clients and processor configurations
	service := processor.New(ctx, apiClients, svcCfg)
	log.Println("processor created")

	router := mux.NewRouter()
	router.HandleFunc("/", jokeHandler(ctx, service)).Methods("GET")

	srv := &http.Server{
		Addr:    ":5000",
		Handler: router,
	}

	go func() {
		log.Println("server created")
		log.Fatal(srv.ListenAndServe())
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-done
}

// Returns an http handler that calls createNamedJoke() and prints the result
func jokeHandler(ctx context.Context, s *processor.JokeService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		joke := createNamedJoke(ctx, s)
		fmt.Fprintf(w, joke)
	}
}

// This function makes a Name Service call, uses the returned name as a request to the Joke Service, and replaces
// `Chuck Norris` with the generated name in the joke
func createNamedJoke(ctx context.Context, service *processor.JokeService) string {
	// Generate name using Name Service
	name, err := service.GenerateName(ctx)
	if err != nil {
		errorType := e.GetType(err)
		switch errorType {
		case e.BadRequest:
			// Return error message but don't kill app
			return "Name API - Bad Request Error\n"
		case e.NotFound:
			return "Name API - Not Found Error\n"
		default:
			// Kill app because error is application level
			log.Fatalf("Could not generate name: ", err)
		}
	}

	// Use name from api call into joke request
	jokeRequest := jokes.JokeReq{
		FirstName: name.Name,
		LastName:  name.Surname,
		LimitTo:   []string{"nerdy"},
	}

	// Generate joke using Joke service
	joke, err := service.GenerateJoke(ctx, jokeRequest)
	if err != nil {
		errorType := e.GetType(err)
		switch errorType {
		case e.BadRequest:
			// Return error message but don't kill app
			return "Joke API - Bad Request Error\n"
		case e.NotFound:
			return "Joke API - Not Found Error\n"
		default:
			// Kill app because error is application level
			log.Fatalf("could not generate joke: %s", err)
		}
	}

	return joke + "\n"
}

func initConfig() (*config.Config, error) {
	var cfg config.Config

	// Normally woukd use os.Getenv to get this, but will just set
	// statically for the scope of this project
	configLocation := fmt.Sprint(".configs/config-e0.yaml")

	// Using viper to read in configuration file which holds base URL
	viper.SetConfigFile(configLocation)
	if err := viper.ReadInConfig(); err != nil {
		return nil, errors.Errorf("failed to read in config: %s", err)
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, errors.Errorf("failed to unmarshal config: %s", err)
	}
	return &cfg, nil
}
