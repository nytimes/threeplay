package threeplay

import (
	"fmt"
	"io/ioutil"
)

// CaptionsFormat is supported output format for captions
type CaptionsFormat string

const (
	// SRT format for captions file
	SRT CaptionsFormat = "srt"
	// DFX format for captions file
	DFX CaptionsFormat = "pdfxp"
	// SMI format for captions file
	SMI CaptionsFormat = "smi"
	// STL format for captions file
	STL CaptionsFormat = "stl"
	// QT format for captions file
	QT CaptionsFormat = "qt"
	// QTXML format for captions file
	QTXML CaptionsFormat = "qtxml"
	// CPTXML format for captions file
	CPTXML CaptionsFormat = "cptxml"
	// ADBE format for captions file
	ADBE CaptionsFormat = "adbe"
)

// GetCaptions get captions by threeplay file ID with specific format
// current supported formats are srt, dfxp, smi, stl, qt, qtxml, cptxml, adbe
func (c *Client) GetCaptions(fileID uint, format CaptionsFormat) ([]byte, error) {
	endpoint := fmt.Sprintf("https://%s/files/%d/captions.%s?apikey=%s",
		threePlayStaticHost, fileID, format, c.apiKey)

	response, err := c.httpClient.Get(endpoint)
	if err != nil {
		return nil, err
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if err := checkForAPIError(responseData); err != nil {
		return nil, err
	}
	return responseData, nil
}

// GetCaptionsByVideoID get captions by video ID with specific format
// current supported formats are srt, dfxp, smi, stl, qt, qtxml, cptxml, adbe
func (c *Client) GetCaptionsByVideoID(id string, format CaptionsFormat) ([]byte, error) {
	endpoint := fmt.Sprintf("https://%s/files/%s/captions.%s?apikey=%s&usevideoid=1",
		threePlayStaticHost, id, format, c.apiKey)

	response, err := c.httpClient.Get(endpoint)
	if err != nil {
		return nil, err
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	if err := checkForAPIError(responseData); err != nil {
		return nil, err
	}
	return responseData, nil
}
