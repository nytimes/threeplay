package v3api_test

import (
	"net/url"
	"testing"

	"github.com/nytimes/threeplay/v3api"
	"github.com/stretchr/testify/assert"
	gock "gopkg.in/h2non/gock.v1"
)

func TestUploadFile(t *testing.T) {
	assert := assert.New(t)

	defer gock.Off()
	gock.New("https://api.3playmedia.com").
		Post("/v3/files").
		MatchType("url").
		BodyString("api_key=api-key&language_id=1&source_id=https%3A%2F%2Fsomewhere.com%2F72397_1_08macron-speech_wg_360p.mp4").
		Reply(200).
		File("../fixtures/v3_file_upload_200.json")

	client := v3api.NewClient("api-key")
	data := url.Values{}
	data.Set("source_id", "https://somewhere.com/72397_1_08macron-speech_wg_360p.mp4")
	data.Set("language_id", "1")

	callParams := v3api.CallParams{APIKey: ""}
	fileID, err := client.UploadFileFromURL(data, callParams)
	assert.Equal(3628518, fileID)
	assert.Nil(err)
}

func TestUploadFileCustomAPI(t *testing.T) {
	assert := assert.New(t)

	defer gock.Off()
	gock.New("https://api.3playmedia.com").
		Post("/v3/files").
		MatchType("url").
		BodyString("api_key=custom-key&language_id=1&source_id=https%3A%2F%2Fsomewhere.com%2F72397_1_08macron-speech_wg_360p.mp4").
		Reply(200).
		File("../fixtures/v3_file_upload_200.json")

	client := v3api.NewClient("api-key")
	data := url.Values{}
	data.Set("source_id", "https://somewhere.com/72397_1_08macron-speech_wg_360p.mp4")
	data.Set("language_id", "1")

	callParams := v3api.CallParams{APIKey: "custom-key"}
	fileID, err := client.UploadFileFromURL(data, callParams)
	assert.Equal(3628518, fileID)
	assert.Nil(err)
}

func TestUploadFileError(t *testing.T) {
	assert := assert.New(t)

	defer gock.Off()
	gock.New("https://api.3playmedia.com").
		Post("/v3/files").
		MatchType("url").
		BodyString("api_key=api-key&bad_param=so-bad&language_id=1").
		Reply(200).
		File("../fixtures/v3_file_upload_400.json")

	client := v3api.NewClient("api-key")
	data := url.Values{}
	data.Set("language_id", "1")
	data.Set("bad_param", "so-bad")

	callParams := v3api.CallParams{APIKey: ""}
	fileID, err := client.UploadFileFromURL(data, callParams)
	assert.Equal(0, fileID)
	assert.NotNil(err)
	assert.Equal("400: parameter_error-Unrecognized parameters supplied: 'bad_param'", err.Error())
}
