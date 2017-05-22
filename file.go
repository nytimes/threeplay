package threeplay

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// File representation
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

// FilesPage representation
type FilesPage struct {
	Files   []File `json:"files"`
	Summary `json:"summary"`
}

// Summary representation
type Summary struct {
	CurrentPage  json.Number `json:"current_page"`
	PerPage      json.Number `json:"per_page"`
	TotalEntries json.Number `json:"total_entries"`
	TotalPages   json.Number `json:"total_pages"`
}

// UpdateFile updates a File metadata
func (c *Client) UpdateFile(fileID uint, data url.Values) error {
	if data == nil {
		return errors.New("Must specify new data")
	}

	data.Set("apikey", c.apiKey)
	data.Set("api_secret_key", c.apiSecret)

	body := strings.NewReader(data.Encode())

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("https://%s/files/%d", threePlayHost, fileID), body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		return err
	}
	response, err := c.httpClient.Do(req)

	responseData, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return err
	}

	// the API returns "1" on success
	if string(responseData) == "1" {
		return nil
	}

	apiError := &Error{}
	err = json.Unmarshal(responseData, apiError)
	if err != nil {
		return err
	}

	if apiError.IsError {
		return errors.New("Api Error")
	}

	return nil
}
