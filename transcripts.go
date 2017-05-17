package threeplay

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
)

// Word words
type Word [2]string

// Transcript transcript
type Transcript struct {
	Words      []Word            `json:"words"`
	Paragraphs []int             `json:"paragraphs"`
	Speakers   map[string]string `json:"speakers"`
}

// GetTranscript get transcript by file ID
func (c *Client) GetTranscript(fileID uint) (*Transcript, error) {
	response, err := c.GetTranscriptWithFormat(fileID, JSON)
	if err != nil {
		return nil, err
	}

	transcript := &Transcript{}
	err = json.Unmarshal(response, transcript)
	if err != nil {
		return nil, err
	}

	return transcript, nil
}

// GetTranscriptWithFormat get transcript by file ID with specific format
func (c *Client) GetTranscriptWithFormat(id uint, format OutputFormat) ([]byte, error) {
	endpoint := fmt.Sprintf("https://%s/files/%d/transcript.%s?apikey=%s",
		threePlayStaticHost, id, format, c.apiKey,
	)

	response, err := c.client.Get(endpoint)
	if err != nil {
		return nil, err
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	apiError := &Error{}
	json.Unmarshal(responseData, apiError)
	if apiError.IsError {
		return nil, errors.New("API Error")
	}

	return responseData, nil
}

// GetTranscriptByVideoID get transcript by video ID
func (c *Client) GetTranscriptByVideoID(videoID string) (*Transcript, error) {

	response, err := c.GetTranscriptByVideoIDWithFormat(videoID, JSON)
	if err != nil {
		return nil, err
	}
	transcript := &Transcript{}
	err = json.Unmarshal(response, transcript)
	if err != nil {
		return nil, err
	}

	return transcript, nil
}

// GetTranscriptByVideoIDWithFormat get transcript by video ID with specific format
func (c *Client) GetTranscriptByVideoIDWithFormat(id string, format OutputFormat) ([]byte, error) {
	endpoint := fmt.Sprintf("https://%s/files/%s/transcript.%s?apikey=%s&usevideoid=1",
		threePlayStaticHost, id, format, c.apiKey,
	)

	response, err := c.client.Get(endpoint)
	if err != nil {
		return nil, err
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	apiError := &Error{}
	json.Unmarshal(responseData, apiError)
	if apiError.IsError {
		return nil, errors.New("API Error")
	}

	return responseData, nil
}
