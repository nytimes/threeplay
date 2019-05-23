package v2 // import "github.com/NYTimes/threeplay/v2

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/nytimes/threeplay/types"
	"github.com/sethgrid/pester"
)

// Client 3Play Media API client
type ClientV2 struct {
	apiKey     string
	apiSecret  string
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
func NewClient(apiKey, apiSecret string) *ClientV2 {
	return NewClientWithHTTPClient(apiKey, apiSecret, &http.Client{Timeout: 10 * time.Second})
}

// NewClientWithHTTPClient returns a 3Play Media client with a custom http client
func NewClientWithHTTPClient(apiKey, apiSecret string, client *http.Client) *ClientV2 {
	httpClient := pester.NewExtendedClient(client)
	httpClient.MaxRetries = 5
	httpClient.Backoff = pester.ExponentialJitterBackoff
	return &ClientV2{
		apiKey:     apiKey,
		apiSecret:  apiSecret,
		httpClient: httpClient,
	}
}

func (c *ClientV2) createURL(endpoint string) url.URL {
	return url.URL{
		Scheme: "https",
		Host:   types.ThreePlayHost,
		Path:   endpoint,
	}
}

func (c *ClientV2) prepareURL(u url.URL, querystrings url.Values) string {
	qs := url.Values{}
	qs.Set("apikey", c.apiKey)
	for k, v := range querystrings {
		qs[k] = v
	}
	u.RawQuery = qs.Encode()
	return u.String()
}

func (c *ClientV2) createRequest(method, endpoint string, data url.Values) (*http.Request, error) {
	data.Set("apikey", c.apiKey)
	data.Set("api_secret_key", c.apiSecret)
	body := strings.NewReader(data.Encode())
	req, err := http.NewRequest(method, endpoint, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req, nil
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
