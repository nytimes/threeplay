package v3api

import (
	"fmt"
	"net/url"
)

type ThreePlayFileResponse struct {
	Code  int                      `json:"code"`
	Data  FileObjectRepresentation `json:"data"`
	Error ThreePlayError           `json:"error"`
}

type FileObjectRepresentation struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Duration    float64 `json:"duration"`
	LanguageID  int     `json:"language_id"`
	LanguageIDs []int   `json:"language_ids"`
	BatchID     int     `json:"batch_id"`
	ReferenceID string  `json:"reference_id"`
}

// UploadFileFromURL uploads a file to threeplay using the file's URL and
// returns the file ID.
func (c *Client) UploadFileFromURL(options url.Values) (int, error) {
	apiURL := c.createURL("/files")
	data := url.Values{}
	data.Set("api_key", c.apiKey)
	for key, val := range options {
		data[key] = val
	}
	res, err := c.httpClient.PostForm(apiURL.String(), data)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	response := &ThreePlayFileResponse{}
	if err := parseResponse(res, response); err != nil {
		return 0, err
	}

	if response.Code != 200 {
		return 0, fmt.Errorf("%v: %v-%v", response.Code, response.Error.Type, response.Error.Message)
	}

	return response.Data.ID, nil
}
