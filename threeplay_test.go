package threeplay

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func createResponseFromJsonFile(jsonFile string) *http.Response {
	file, _ := ioutil.ReadFile(jsonFile)
	data := bytes.NewReader(file)
	resp := http.Response{Body: ioutil.NopCloser(data)}
	return &resp
}

type HTTPClientMock struct {
	mock.Mock
}

func (c *HTTPClientMock) Get(url string) (*http.Response, error) {
	args := c.Called(url)
	return args.Get(0).(*http.Response), nil
}

func (c *HTTPClientMock) PostForm(url string, data url.Values) (*http.Response, error) {
	args := c.Called(url, data)
	return args.Get(0).(*http.Response), nil
}

func TestGetFile(t *testing.T) {
	httpClient := &HTTPClientMock{}
	expectedApiCall := "https://api.3playmedia.com/files/123456?apikey=api-key"
	c.On("Get", expectedApiCall).Return(createResponseFromJsonFile("./fixtures/file.json"), nil)
	client := NewClientWithHTTPClient("api-key", "", c)
	file, _ := client.GetFile(123456)
	assert.Equal(t, file.Name, "72397_1_08macron-speech_wg_360p.mp4")
	httpClient.AssertExpectations(t)
}

func TestGetFiles(t *testing.T) {
	httpClient := &HTTPClientMock{}
	expectedApiCall := "https://api.3playmedia.com/files?apikey=api-key"
	httpClient.On("Get", expectedApiCall).Return(createResponseFromJsonFile("./fixtures/files_page1.json"), nil)
	client := NewClientWithHTTPClient("api-key", httpClient)

	filesPage, _ := client.GetFiles()

	assert.Equal(t, len(filesPage.Files), 10)
	assert.Equal(t, filesPage.Summary.CurrentPage.String(), "1")
	assert.Equal(t, filesPage.Summary.PerPage.String(), "10")
	httpClient.AssertExpectations(t)
}

func TestGetFilesWithPagination(t *testing.T) {
	httpClient := &HTTPClientMock{}
	expectedApiCall := "https://api.3playmedia.com/files?apikey=api-key&page=2"
	httpClient.On("Get", expectedApiCall).Return(createResponseFromJsonFile("./fixtures/files_page2.json"), nil)
	client := NewClientWithHTTPClient("api-key", httpClient)
	querystring := url.Values{}
	querystring.Add("page", "2")
	filesPage, _ := client.GetFiles(querystring)
	assert.Equal(t, "2", filesPage.Summary.CurrentPage.String())
	httpClient.AssertExpectations(t)
}

func TestUploadFileFromURL(t *testing.T) {
	assert := assert.New(t)
	c := &HTTPClientMock{}

	expectedEndpoint := "https://api.3playmedia.com/files"
	expectedData := url.Values{}
	expectedData.Set("apikey", ":api-key")
	expectedData.Set("api_secret_key", ":secret")
	expectedData.Set("link", "https://somewhere.com/72397_1_08macron-speech_wg_360p.mp4")
	expectedData.Set("video_id", "123456")

	fakeResponse := &http.Response{
		Body: ioutil.NopCloser(bytes.NewBufferString("1686514")),
	}
	c.On("PostForm", expectedEndpoint, expectedData).Return(fakeResponse, nil)
	client := NewClientWithHTTPClient(":api-key", ":secret", c)
	data := url.Values{}
	data.Set("video_id", "123456")

	fileID, _ := client.UploadFileFromURL("https://somewhere.com/72397_1_08macron-speech_wg_360p.mp4", data)
	assert.Equal("1686514", fileID)
}
