package threeplay

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
)

// GetTags gets the list of tags of a file
func (c *Client) GetTags(fileID uint) ([]string, error) {

	endpoint := fmt.Sprintf("https://%s/files/%d/tags?apikey=%s", threePlayHost, fileID, c.apiKey)
	response, err := c.httpClient.Get(endpoint)
	if err != nil {
		return nil, err
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	apiError := &Error{}
	err = json.Unmarshal(responseData, apiError)
	if err == nil && apiError.IsError {
		return nil, errors.New("Api Error")
	}

	var tags []string
	err = json.Unmarshal(responseData, &tags)
	if err != nil {
		return nil, err
	}

	return tags, nil
}

type addTagResult struct {
	Result bool     `json:"result"`
	Tags   []string `json:"media_file_tags"`
}

// AddTag adds a tag to a file
func (c *Client) AddTag(fileID uint, tag string) ([]string, error) {
	endpoint := fmt.Sprintf("https://%s/files/%d/tags", threePlayHost, fileID)

	data := url.Values{}
	data.Set("apikey", c.apiKey)
	data.Set("api_secret_key", c.apiSecret)
	data.Set("name", tag)

	response, err := c.httpClient.PostForm(endpoint, data)
	if err != nil {
		return nil, err
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	apiError := &Error{}
	err = json.Unmarshal(responseData, apiError)
	if err == nil && apiError.IsError {
		return nil, errors.New("Api Error")
	}

	result := &addTagResult{}
	err = json.Unmarshal(responseData, &result)
	if err != nil {
		return nil, err
	}

	if result.Result != true {
		return nil, errors.New("Adding Tag Failed")
	}

	return result.Tags, nil
}

// RemoveTag removes a tag of a file
func (c *Client) RemoveTag(fileID uint, tag string) ([]string, error) {
	endpoint := fmt.Sprintf("https://%s/files/%d/tags/%s", threePlayHost, fileID, tag)

	data := url.Values{}
	data.Set("apikey", c.apiKey)
	data.Set("api_secret_key", c.apiSecret)
	data.Set("_method", "delete")

	response, err := c.httpClient.PostForm(endpoint, data)
	if err != nil {
		return nil, err
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	apiError := &Error{}
	err = json.Unmarshal(responseData, apiError)
	if err == nil && apiError.IsError {
		return nil, errors.New("Api Error")
	}

	var tags []string
	err = json.Unmarshal(responseData, &tags)
	if err != nil {
		return nil, err
	}

	return tags, nil
}
