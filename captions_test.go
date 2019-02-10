package threeplay_test

import (
	"testing"

	"github.com/NYTimes/threeplay"
	"github.com/stretchr/testify/assert"
	gock "gopkg.in/h2non/gock.v1"
)

func TestGetCaptions(t *testing.T) {
	var tests = []struct {
		name   string
		opts   threeplay.GetCaptionsOptions
		path   string
		params map[string]string
	}{
		{
			"with file id and standard format",
			threeplay.GetCaptionsOptions{
				FileID: 123456,
				Format: threeplay.SRT,
			},
			"/files/123456/captions.srt",
			map[string]string{"apikey": "api-key"},
		},
		{
			"with file id and custom format",
			threeplay.GetCaptionsOptions{
				FileID:       123456,
				OutputFormat: "42.srt",
			},
			"/files/123456/output_formats/42.srt",
			map[string]string{"apikey": "api-key"},
		},
		{
			"with video id and standard format",
			threeplay.GetCaptionsOptions{
				VideoID: "vid-123",
				Format:  threeplay.SRT,
			},
			"/files/vid-123/captions.srt",
			map[string]string{"apikey": "api-key", "usevideoid": "1"},
		},
		{
			"with video id and custom format",
			threeplay.GetCaptionsOptions{
				VideoID:      "vid-123",
				OutputFormat: "42.vtt",
			},
			"/files/vid-123/output_formats/42.vtt",
			map[string]string{"apikey": "api-key", "usevideoid": "1"},
		},
		{
			"with all fields - should prefer custom format && file id",
			threeplay.GetCaptionsOptions{
				FileID:       123456,
				VideoID:      "vid-123",
				Format:       threeplay.WebVTT,
				OutputFormat: "42.srt",
			},
			"/files/123456/output_formats/42.srt",
			map[string]string{"apikey": "api-key"},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)
			defer gock.Off()

			gock.New("https://static.3playmedia.com").
				Get(test.path).
				MatchParams(test.params).
				Reply(200).
				File("./fixtures/captions.srt")

			client := threeplay.NewClient("api-key", "secret-key")
			result, err := client.GetCaptions(test.opts)
			assert.NotNil(result)
			assert.Nil(err)
		})
	}
}

func TestGetCaptionsApiInvalidOptions(t *testing.T) {
	var tests = []struct {
		name string
		opts threeplay.GetCaptionsOptions
	}{
		{
			"missing format",
			threeplay.GetCaptionsOptions{
				FileID:  10,
				VideoID: "vid-123",
			},
		},
		{
			"missing id",
			threeplay.GetCaptionsOptions{
				Format:       threeplay.SRT,
				OutputFormat: "42.srt",
			},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)
			client := threeplay.NewClient("api-key", "secret-key")
			result, err := client.GetCaptions(test.opts)
			assert.Nil(result)
			assert.NotNil(err)
		})
	}
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
	result, err := client.GetCaptions(threeplay.GetCaptionsOptions{
		FileID: 123456,
		Format: threeplay.SRT,
	})
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
