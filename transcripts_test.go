package threeplay_test

import (
	"fmt"
	"testing"

	"github.com/NYTimes/threeplay"
	"github.com/stretchr/testify/assert"
	gock "gopkg.in/h2non/gock.v1"
)

func TestGetTranscriptWithFormat(t *testing.T) {
	assert := assert.New(t)
	defer gock.Off()

	expectedResult := "some-transcript-data"

	gock.New("https://static.3playmedia.com").
		Get("/files/123456/transcript.txt").
		MatchParam("apikey", "api-key").
		Reply(200).
		BodyString(expectedResult)

	client := threeplay.NewClient("api-key", "secret-key")
	result, err := client.GetTranscriptWithFormat(123456, threeplay.TXT)
	assert.Equal(expectedResult, string(result))
	assert.Nil(err)
}

func TestGetTranscriptWithFormatApiError(t *testing.T) {
	assert := assert.New(t)
	defer gock.Off()

	gock.New("https://static.3playmedia.com").
		Get("/files/123456/transcript.txt").
		MatchParam("apikey", "api-key").
		Reply(200).
		File("./fixtures/error.json")

	client := threeplay.NewClient("api-key", "secret-key")
	result, err := client.GetTranscriptWithFormat(123456, threeplay.TXT)
	assert.Nil(result)
	assert.NotNil(err)
	assert.Equal(threeplay.ErrUnauthorized.Error(), err.Error())
}

func TestGetTranscript(t *testing.T) {
	assert := assert.New(t)
	defer gock.Off()

	gock.New("https://static.3playmedia.com").
		Get("/files/123456/transcript.json").
		MatchParam("apikey", "api-key").
		Reply(200).
		File("./fixtures/transcript.json")

	client := threeplay.NewClient("api-key", "secret-key")
	transcript, err := client.GetTranscript(123456)
	assert.NotNil(transcript)
	assert.Nil(err)
}

func TestGetTranscriptApiError(t *testing.T) {
	assert := assert.New(t)
	defer gock.Off()

	gock.New("https://static.3playmedia.com").
		Get("/files/123456/transcript.json").
		MatchParam("apikey", "api-key").
		Reply(200).
		File("./fixtures/error.json")

	client := threeplay.NewClient("api-key", "secret-key")
	transcript, err := client.GetTranscript(123456)
	assert.Nil(transcript)
	assert.NotNil(err)
	assert.Equal(threeplay.ErrUnauthorized.Error(), err.Error())
}

func TestGetTranscriptByVideoID(t *testing.T) {
	assert := assert.New(t)
	defer gock.Off()

	gock.New("https://static.3playmedia.com").
		Get("/files/123456/transcript.json").
		MatchParam("apikey", "api-key").
		MatchParam("usevideoid", "1").
		Reply(200).
		File("./fixtures/transcript.json")

	client := threeplay.NewClient("api-key", "secret-key")
	transcript, err := client.GetTranscriptByVideoID("123456")
	assert.NotNil(transcript)
	assert.Nil(err)
}

func TestGetTranscriptByVideoIDApiError(t *testing.T) {
	assert := assert.New(t)
	defer gock.Off()

	gock.New("https://static.3playmedia.com").
		Get("/files/123456/transcript.json").
		MatchParam("apikey", "api-key").
		MatchParam("usevideoid", "1").
		Reply(200).
		File("./fixtures/error.json")

	client := threeplay.NewClient("api-key", "secret-key")
	transcript, err := client.GetTranscriptByVideoID("123456")
	assert.Nil(transcript)
	assert.NotNil(err)
	assert.Equal(threeplay.ErrUnauthorized.Error(), err.Error())
}

func TestGetTranscriptByVideoIDWithFormat(t *testing.T) {
	assert := assert.New(t)
	defer gock.Off()

	expectedResult := "some-transcript-data"

	gock.New("https://static.3playmedia.com").
		Get("/files/123456/transcript.txt").
		MatchParam("apikey", "api-key").
		MatchParam("usevideoid", "1").
		Reply(200).
		BodyString(expectedResult)

	client := threeplay.NewClient("api-key", "secret-key")
	result, err := client.GetTranscriptByVideoIDWithFormat("123456", threeplay.TXT)
	assert.Equal(expectedResult, string(result))
	assert.Nil(err)
}

func TestGetTranscriptByVideoIDtWithFormatApiError(t *testing.T) {
	assert := assert.New(t)
	defer gock.Off()

	gock.New("https://static.3playmedia.com").
		Get("/files/123456/transcript.txt").
		MatchParam("apikey", "api-key").
		MatchParam("usevideoid", "1").
		Reply(200).
		File("./fixtures/error.json")

	client := threeplay.NewClient("api-key", "secret-key")
	result, err := client.GetTranscriptByVideoIDWithFormat("123456", threeplay.TXT)
	assert.Nil(result)
	assert.NotNil(err)
	assert.Equal(threeplay.ErrUnauthorized.Error(), err.Error())
}

func ExampleClient_GetTranscriptWithFormat() {
	client := threeplay.NewClient("api-key", "secret")
	transcript, _ := client.GetTranscriptWithFormat(123, threeplay.JSON)
	fmt.Println(transcript)
}
