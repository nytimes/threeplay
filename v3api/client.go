package v3api // import "github.com/NYTimes/threeplay

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/nytimes/threeplay/types"
	"github.com/sethgrid/pester"
)

// Client 3Play Media API client
type Client struct {
	apiKey     string
	httpClient *pester.Client
}

// Error representation of 3Play API error
type Error struct {
	IsError bool              `json:"iserror"`
	Errors  map[string]string `json:"errors"`
}

var (
	// ErrUnauthorized represents a 401 error on API
	ErrUnauthorized = errors.New("401: API Error")
	// ErrNotFound represents a 404 error on API
	ErrNotFound = errors.New("404: API Error")
)

// NewClient returns a 3Play Media client
func NewClient(apiKey string) *Client {
	return NewClientWithHTTPClient(apiKey, &http.Client{Timeout: 10 * time.Second})
}

// NewClientWithHTTPClient returns a 3Play Media client with a custom http client
func NewClientWithHTTPClient(apiKey string, client *http.Client) *Client {
	httpClient := pester.NewExtendedClient(client)
	httpClient.MaxRetries = 5
	httpClient.Backoff = pester.ExponentialJitterBackoff
	return &Client{
		apiKey:     apiKey,
		httpClient: httpClient,
	}
}

func (c *Client) createURL(endpoint string) url.URL {
	return url.URL{
		Scheme: "https",
		Host:   types.ThreePlayHost,
		Path:   fmt.Sprintf("/v3%v", endpoint),
	}
}

func parseResponse(res *http.Response, ref interface{}) error {
	responseData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	err = checkForAPIError(responseData)
	if err != nil {
		return err
	}
	err = json.Unmarshal(responseData, ref)

	return err
}

func checkForAPIError(responseData []byte) error {
	apiError := Error{}
	json.Unmarshal(responseData, &apiError)
	if apiError.IsError {
		if _, ok := apiError.Errors["authentication"]; ok {
			return ErrUnauthorized
		}
		if _, ok := apiError.Errors["not_found"]; ok {
			return ErrNotFound
		}

		return errors.New("api error: " + string(responseData))
	}
	return nil
}
