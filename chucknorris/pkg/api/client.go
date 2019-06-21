package api

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	e "chucknorris/pkg/errors"

	"github.com/jpillora/backoff"
	"github.com/pkg/errors"
)

// Client holds the request
type Client interface {
	Request(ctx context.Context, req Request, response interface{}) (err error)
}

// Config holds the API configuration data
type Config struct {
	Name    string `yaml:"name" json:"name"`
	BaseURL string `yaml:"base_url" json:"base_url"`
	Timeout int    `yaml:"timeout" json:"timeout"`
}

// Request holds the http request fiels
type Request struct {
	Method, URL string
	Query       map[string][]string
	Body        io.Reader
}

// RstService holds the configuration and http client data
type RstService struct {
	cfg       *Config
	rstClient *http.Client
}

// New returns a new RstService
func (c *Config) New() *RstService {
	return &RstService{
		cfg: c,
		// Set timeout to 5 seconds, using config .yaml file
		rstClient: &http.Client{
			Timeout: time.Duration(c.Timeout) * time.Millisecond},
	}
}

// Request generates the request that will be executed
func (rstSvc *RstService) Request(ctx context.Context, req Request, response interface{}) (err error) {
	var resErr error

	// Request will be retried 2 times if request response code is in 400 range
	for retry := 0; retry < 3; retry++ {
		if req.Body != nil {
			body, err := ioutil.ReadAll(req.Body)
			if err != nil {
				// Wrap error and restart loop on next iteration
				return e.NewfErr("unable to read body: %s", err.Error())
			}
			req.Body = bytes.NewReader(body)
		}

		log.Println("performing HTTP request")

		// create url
		u, err := url.Parse(req.URL)
		if err != nil {
			return e.NewfErr("unable to parse URL: %s", err)
		}

		// add query params to url
		q := u.Query()
		if req.Query != nil {
			for k, v := range req.Query {
				for _, vv := range v {
					q.Add(k, vv)
				}
			}
			u.RawQuery = q.Encode()
		}

		// create request for service call
		var r *http.Request
		r, err = http.NewRequest(req.Method, u.String(), req.Body)
		if err != nil {
			return e.NewfErr("unable to create request: %s", err.Error())
		}

		// call the service
		res, err := rstSvc.rstClient.Do(r)
		if err != nil {
			return e.NewfErr("unable to complete request: %s", err.Error())
		}
		defer res.Body.Close()

		// Imported library for exponential backoff, will start at 1 second, and double in size
		// until it reaches a max of 10 seconds
		b := &backoff.Backoff{
			Min:    1000 * time.Millisecond,
			Max:    10 * time.Second,
			Factor: 2,
			Jitter: false,
		}

		// backoff duration
		d := b.Duration()

		var errType e.ErrorType

		// Ensure status code was 200, else continue into next iteration of retry for-loop
		switch res.StatusCode / 100 {
		case 2: // 2** success cases are ok
			// reset the exponential backoff counter and return response
			b.Reset()
			return json.NewDecoder(res.Body).Decode(response)
		case 3: // 3** redirect cases are ok, could be a redirect
			b.Reset()
			return json.NewDecoder(res.Body).Decode(response)
		case 4: // 4** client side errors are not okay, retryable
			// Use backoff to determine how long sleep duration between retries is
			time.Sleep(d)
			// Get body response for errors
			bdy, err := bodyString(res)
			if err != nil {
				return e.NewErr("unable to read body of response")
			}
			errType = e.BadRequest
			resErr = errType.New("Bad Request: "+bdy, res.StatusCode)
			log.Printf("Bad Request: %s, Status Code: %d", bdy, res.StatusCode)
		case 5: // 5** server side errors are not okay, not retryable
			time.Sleep(d)
			bdy, err := bodyString(res)
			if err != nil {
				return e.NewErr("unable to read body of response")
			}
			errType = e.NotFound
			resErr = errType.New("Not Found: "+bdy, res.StatusCode)
			log.Printf("Not Found: %s, Status Code: %d", bdy, res.StatusCode)
		default: // Unknown response code, retry
			time.Sleep(d)
			bdy, err := bodyString(res)
			if err != nil {
				return e.NewErr("unable to read body of response")
			}
			errType = e.NotFound
			resErr = errType.New("Unknown Response: "+bdy, res.StatusCode)
			log.Printf("Unknown Response: %s, Status Code: %d", bdy, res.StatusCode)
		}
		// Reset exponential backoff when loop completes
		b.Reset()
	}
	// return error if one exists
	return resErr
}

// Returns the body of a response as a string
func bodyString(res *http.Response) (string, error) {
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", errors.Wrap(err, "unable to read body")
	}
	bodyString := string(bodyBytes)

	return bodyString, nil
}
