package threeplay

import (
	"bytes"
	"io/ioutil"
	"net/http"
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

func TestGetFile(t *testing.T) {
	c := &HTTPClientMock{}
	expectedApiCall := "https://api.3playmedia.com/files/123456?apikey=api-key"
	c.On("Get", expectedApiCall).Return(createResponseFromJsonFile("./fixtures/file.json"), nil)
	client := NewClientWithHTTPClient("api-key", c)
	file, _ := client.GetFile(123456)
	assert.Equal(t, file.Name, "72397_1_08macron-speech_wg_360p.mp4")
	c.AssertExpectations(t)
}
