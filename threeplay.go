package threeplay

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
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
	Get(string) (resp *http.Response, err error)
}

type Client struct {
	apiKey string
	client HTTPClient
}

func NewClient(apiKey string) Client {
	return Client{
		apiKey: apiKey,
		client: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

func NewClientWithHTTPClient(apiKey string, client HTTPClient) Client {
	return Client{
		apiKey: apiKey,
		client: client,
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
