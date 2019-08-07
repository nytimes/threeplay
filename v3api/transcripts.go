package v3api

import (
	"fmt"
	"net/url"

	"github.com/nytimes/threeplay/types"
)

// TranscriptFormatToID maps a caption format to their 3play code
var TranscriptFormatToID = map[types.CaptionsFormat]int{
	types.WebVTT: 139,
	types.SRT:    7,
}

// ThreePlayTranscriptResponse the info response object
type ThreePlayTranscriptResponse struct {
	Code  int                            `json:"code"`
	Data  TranscriptObjectRepresentation `json:"data"`
	Error ThreePlayError                 `json:"error"`
}

// ThreePlayTranscriptTextResponse is the text of the transcript response
type ThreePlayTranscriptTextResponse struct {
	Code  int            `json:"code"`
	Data  string         `json:"data"`
	Error ThreePlayError `json:"error"`
}

// ThreePlayTranscriptCancelResponse the cancel response object
type ThreePlayTranscriptCancelResponse struct {
	Code  int                        `json:"code"`
	Data  CancelObjectRepresentation `json:"data"`
	Error ThreePlayError             `json:"error"`
}

// ThreePlayError represents the content of a transcript error response
type ThreePlayError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

// TranscriptObjectRepresentation represents the content of a transcript info response
type TranscriptObjectRepresentation struct {
	ID                  int     `json:"id"`
	MediaFileID         int     `json:"media_file_id"`
	Duration            float64 `json:"duration"`
	Type                string  `json:"type"`
	LanguageID          int     `json:"language_id"`
	Status              string  `json:"status"`
	Cancellable         bool    `json:"cancellable"`
	CancellationReason  string  `json:"cancellation_reason"`
	CancellationDetails string  `json:"cancellation_details"`
}

// CancelObjectRepresentation represents the content of a cancel response
type CancelObjectRepresentation struct {
	Success bool `json:"success"`
}

// CallParams contains the call params the caller wants to override
type CallParams struct {
	APIKey string `json:"api_key"`
}

// OrderTranscript orders a transcript generation job
func (c *Client) OrderTranscript(mediaFileID, callbackURL, turnaroundLevel string, callParams CallParams) (*TranscriptObjectRepresentation, error) {
	var apiURL url.URL
	apiKey := c.setAPIKey(callParams.APIKey)
	data := url.Values{}

	data.Set("api_key", apiKey)
	data.Set("media_file_id", mediaFileID)
	if len(callbackURL) > 0 {
		data.Set("callback", callbackURL)
	}
	if turnaroundLevel == "asr" {
		apiURL = c.createURL("/transcripts/order/asr")
	} else {
		apiURL = c.createURL("/transcripts/order/transcription")
		data.Set("turnaround_level_id", turnaroundLevel)
	}
	res, err := c.httpClient.PostForm(apiURL.String(), data)
	if err != nil {
		return &TranscriptObjectRepresentation{}, err
	}
	response := &ThreePlayTranscriptResponse{}
	if err := parseResponse(res, response); err != nil {
		return &TranscriptObjectRepresentation{}, err
	}
	if response.Code != 200 {
		return &TranscriptObjectRepresentation{}, fmt.Errorf("%v: %v-%v", response.Code, response.Error.Type, response.Error.Message)
	}
	return &response.Data, nil
}

// GetTranscriptInfo gets the status of the transcript job
func (c *Client) GetTranscriptInfo(mediaFileID string, callParams CallParams) (*TranscriptObjectRepresentation, error) {
	apiKey := c.setAPIKey(callParams.APIKey)
	endpoint := fmt.Sprintf("https://%s/v3/transcripts/%s?api_key=%s",
		types.ThreePlayHost, mediaFileID, apiKey,
	)

	res, err := c.httpClient.Get(endpoint)
	if err != nil {
		return nil, err
	}
	response := &ThreePlayTranscriptResponse{}
	if err := parseResponse(res, response); err != nil {
		return &TranscriptObjectRepresentation{}, err
	}
	if response.Code != 200 {
		return &TranscriptObjectRepresentation{}, fmt.Errorf("%v: %v-%v", response.Code, response.Error.Type, response.Error.Message)
	}
	return &response.Data, nil
}

// GetTranscriptText downloads the transcript in the specified format
func (c *Client) GetTranscriptText(mediaFileID, offset string, outputFormat types.CaptionsFormat, callParams CallParams) (string, error) {
	apiKey := c.setAPIKey(callParams.APIKey)
	endpoint := fmt.Sprintf("https://%s/v3/transcripts/%s/text?api_key=%s&output_format_id=%d",
		types.ThreePlayHost, mediaFileID, apiKey, TranscriptFormatToID[outputFormat],
	)
	if offset != "" {
		endpoint += fmt.Sprintf("&offset=%s", offset)
	}
	res, err := c.httpClient.Get(endpoint)
	if err != nil {
		return "", err
	}

	response := &ThreePlayTranscriptTextResponse{}
	if err := parseResponse(res, response); err != nil {
		return "", err
	}
	if response.Code != 200 {
		return "", fmt.Errorf("%v: %v-%v", response.Code, response.Error.Type, response.Error.Message)
	}
	return response.Data, nil
}

// CancelTranscript cancels the transcript order if possible
func (c *Client) CancelTranscript(mediaFileID string, callParams CallParams) error {
	apiKey := c.setAPIKey(callParams.APIKey)
	apiURL := c.createURL(fmt.Sprintf("/transcripts/%s/cancel", mediaFileID))
	data := url.Values{}
	data.Set("api_key", apiKey)
	res, err := c.httpClient.PostForm(apiURL.String(), data)
	if err != nil {
		return err
	}
	response := &ThreePlayTranscriptCancelResponse{}
	if err := parseResponse(res, response); err != nil {
		return err
	}
	if response.Code != 200 {
		return fmt.Errorf("%v: %v-%v", response.Code, response.Error.Type, response.Error.Message)
	}
	return nil
}

// GetEditingLink gets an expiring editing link
func (c *Client) GetEditingLink(mediaFileID string, hoursUntilExpiration int, callParams CallParams) (string, error) {
	apiKey := c.setAPIKey(callParams.APIKey)
	endpoint := fmt.Sprintf("https://%s/v3/transcripts/%s/expiring_editing_link?api_key=%s&hours_until_expiration=%d",
		types.ThreePlayHost, mediaFileID, apiKey, hoursUntilExpiration,
	)
	res, err := c.httpClient.Get(endpoint)
	if err != nil {
		return "", err
	}
	response := &ThreePlayTranscriptTextResponse{}
	if err := parseResponse(res, response); err != nil {
		return "", nil
	}
	if response.Code != 200 {
		return "", fmt.Errorf("%v: %v-%v", response.Code, response.Error.Type, response.Error.Message)
	}
	return response.Data, nil
}

func (c *Client) setAPIKey(apiKeyParam string) string {
	if len(apiKeyParam) > 0 {
		return apiKeyParam
	}
	return c.apiKey
}
