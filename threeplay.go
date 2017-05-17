package threeplay

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const threePlayHost = "api.3playmedia.com"
const threePlayStaticHost = "static.3playmedia.com"

// OutputFormat supported output formats for transcriptions
type OutputFormat string

const (
	// JSON format for transcripted file
	JSON OutputFormat = "json"
	// TXT format output for transcripted file
	TXT OutputFormat = "txt"
	// HTML format output for transcripted file
	HTML OutputFormat = "html"
)

// HTTPClient interface for http clients used in the 3Play client
type HTTPClient interface {
	Get(string) (*http.Response, error)
	PostForm(string, url.Values) (*http.Response, error)
}

// Client 3Play Media API client
type Client struct {
	apiKey    string
	apiSecret string
	client    HTTPClient
}

// Error representation of 3Play API error
type Error struct {
	IsError bool              `json:"iserror"`
	Errors  map[string]string `json:"errors"`
}

// NewClient returns a 3Play Media client
func NewClient(apiKey, apiSecret string) *Client {
	return &Client{
		apiKey:    apiKey,
		apiSecret: apiSecret,
		client: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

// NewClientWithHTTPClient returns a 3Play Media client with a custom http client
func NewClientWithHTTPClient(apiKey, apiSecret string, client HTTPClient) *Client {
	return &Client{
		apiKey:    apiKey,
		apiSecret: apiSecret,
		client:    client,
	}
}

func (c Client) buildURL(endpoint string, querystring url.Values) string {
	querystring.Add("apikey", c.apiKey)

	url := url.URL{
		Scheme:   "https",
		Host:     threePlayHost,
		Path:     endpoint,
		RawQuery: querystring.Encode(),
	}

	return url.String()
}

func (c Client) fetchAndParse(endpoint string, ref interface{}) error {
	apiError := &Error{}
	response, err := c.client.Get(endpoint)

	if err != nil {
		return err
	}

	responseData, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return err
	}

	err = json.Unmarshal(responseData, apiError)

	if err != nil {
		return err
	}

	if apiError.IsError {
		return errors.New("API Error")
	}

	err = json.Unmarshal(responseData, ref)

	if err != nil {
		return err
	}

	return nil
}

// GetFiles returns a list of files
func (c *Client) GetFiles(params url.Values) (*FilesPage, error) {
	querystring := url.Values{}
	if params != nil {
		querystring = params
	}

	filesPage := &FilesPage{}
	endpoint := c.buildURL("/files", querystring)
	if err := c.fetchAndParse(endpoint, filesPage); err != nil {
		return nil, err
	}
	return filesPage, nil
}

// GetFile gets a single file by `file_id`
func (c *Client) GetFile(id uint) (*File, error) {
	file := &File{}
	endpoint := c.buildURL(fmt.Sprintf("/files/%d", id), url.Values{})
	if err := c.fetchAndParse(endpoint, file); err != nil {
		return nil, err
	}
	return file, nil
}

// UploadFileFromURL uploads a file to threeplay using the file's URL.
func (c *Client) UploadFileFromURL(fileURL string, options url.Values) (string, error) {
	endpoint := fmt.Sprintf("https://%s/files", threePlayHost)

	data := url.Values{}
	data.Set("apikey", c.apiKey)
	data.Set("api_secret_key", c.apiSecret)
	data.Set("link", fileURL)

	for key, val := range options {
		data[key] = val
	}

	response, err := c.client.PostForm(endpoint, data)
	if err != nil {
		return "", err
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	apiError := &Error{}
	json.Unmarshal(responseData, apiError)
	if apiError.IsError {
		return "", errors.New("API Error")
	}

	return string(responseData), nil
}
