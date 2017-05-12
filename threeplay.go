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

type File struct {
	ID                   uint   `json:"id"`
	ProjectID            uint   `json:"project_id"`
	BatchID              uint   `json:"batch_id"`
	Duration             uint   `json:"duration"`
	Attribute1           string `json:"attribute1"`
	Attribute2           string `json:"attribute2"`
	Attribute3           string `json:"attribute3"`
	VideoID              string `json:"video_id"`
	Name                 string `json:"name"`
	CallbackURL          string `json:"callback_url"`
	Description          string `json:"description"`
	CreatedAt            string `json:"created_at"`
	UpdatedAt            string `json:"updated_at"`
	WordCount            uint   `json:"word_count"`
	ThumbnailURL         string `json:"thumbnail_url"`
	LanguageID           int    `json:"language_id"`
	DefaultServiceTypeID int    `json:"default_service_type_id"`
	Downloaded           bool   `json:"downloaded"`
	State                string `json:"state"`
	TurnaroundLevelID    int    `json:"turnaround_level_id"`
	Deadline             string `json:"deadline"`
	BatchName            string `json:"batch_name"`
	ErrorDescription     string `json:"error_description"`
}

type HTTPClient interface {
	Get(string) (*http.Response, error)
	PostForm(string, url.Values) (*http.Response, error)
}

type Client struct {
	apiKey    string
	apiSecret string
	client    HTTPClient
}

func NewClient(apiKey, apiSecret string) Client {
	return Client{
		apiKey:    apiKey,
		apiSecret: apiSecret,
		client: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

func NewClientWithHTTPClient(apiKey, apiSecret string, client HTTPClient) Client {
	return Client{
		apiKey:    apiKey,
		apiSecret: apiSecret,
		client:    client,
	}
}

const ThreePlayHost = "https://api.3playmedia.com"

type Error struct {
	IsError bool              `json:"iserror"`
	Errors  map[string]string `json:"errors"`
}

func (c Client) GetFile(id uint) (*File, error) {
	file := &File{}
	apiError := &Error{}
	endpoint := fmt.Sprintf("%s/files/%d?apikey=%s", ThreePlayHost, id, c.apiKey)
	response, err := c.client.Get(endpoint)

	if err != nil {
		return nil, err
	}

	responseData, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(responseData, apiError)

	if err != nil {
		return nil, err
	}

	if apiError.IsError {
		return nil, errors.New("API Error")
	}

	err = json.Unmarshal(responseData, file)

	if err != nil {
		return nil, err
	}

	return file, nil
}

//UploadFile uploads a file to threeplay using the file's URL.
func (c *Client) UploadFile(fileURL string, options url.Values) (string, error) {
	endpoint := fmt.Sprintf("%s/files", ThreePlayHost)

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
