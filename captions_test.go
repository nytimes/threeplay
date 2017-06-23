package threeplay_test

import (
	"testing"

	"github.com/NYTimes/threeplay"
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

func TestGetCaptions(t *testing.T) {
	assert := assert.New(t)
	defer gock.Off()

	gock.New("https://static.3playmedia.com").
		Get("/files/123456/captions.srt").
		MatchParam("apikey", "api-key").
		Reply(200).
		File("./fixtures/captions.srt")

	client := threeplay.NewClient("api-key", "secret-key")
	result, err := client.GetCaptions(123456, threeplay.SRT)
	assert.NotNil(result)
	assert.Nil(err)
}

func TestGetCaptionsApiError(t *testing.T) {
	assert := assert.New(t)
	defer gock.Off()

	gock.New("https://static.3playmedia.com").
		Get("/files/123456/captions.srt").
		MatchParam("apikey", "api-key").
		Reply(200).
		File("./fixtures/error.json")

	client := threeplay.NewClient("api-key", "secret-key")
	result, err := client.GetCaptions(123456, threeplay.SRT)
	assert.Nil(result)
	assert.NotNil(err)
	assert.Equal(threeplay.ErrUnauthorized.Error(), err.Error())
}

func TestGetCaptionsByVideoID(t *testing.T) {
	assert := assert.New(t)
	defer gock.Off()

	gock.New("https://static.3playmedia.com").
		Get("/files/123456/captions.srt").
		MatchParam("apikey", "api-key").
		MatchParam("usevideoid", "1").
		Reply(200).
		File("./fixtures/captions.srt")

	client := threeplay.NewClient("api-key", "secret-key")
	result, err := client.GetCaptionsByVideoID("123456", threeplay.SRT)
	assert.NotNil(result)
	assert.Nil(err)
}

func TestGetCaptionsByVideoIDApiError(t *testing.T) {
	assert := assert.New(t)
	defer gock.Off()

	gock.New("https://static.3playmedia.com").
		Get("/files/123456/captions.srt").
		MatchParam("apikey", "api-key").
		MatchParam("usevideoid", "1").
		Reply(200).
		File("./fixtures/error.json")

	client := threeplay.NewClient("api-key", "secret-key")
	result, err := client.GetCaptionsByVideoID("123456", threeplay.SRT)
	assert.Nil(result)
	assert.NotNil(err)
	assert.Equal(threeplay.ErrUnauthorized.Error(), err.Error())
}
