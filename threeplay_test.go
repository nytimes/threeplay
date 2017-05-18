package threeplay_test

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/NYTimes/threeplay"
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

func TestGetFile(t *testing.T) {
	assert := assert.New(t)

	defer gock.Off()
	gock.New("https://api.3playmedia.com").
		Get("/files/123456").
		MatchParam("apikey", "api-key").
		Reply(200).
		File("./fixtures/file.json")

	client := threeplay.NewClient("api-key", "secret-key")
	file, err := client.GetFile(123456)
	assert.Equal(file.Name, "72397_1_08macron-speech_wg_360p.mp4")
	assert.Nil(err)
}

func TestGetFileAPIError(t *testing.T) {
	assert := assert.New(t)

	defer gock.Off()
	gock.New("https://api.3playmedia.com").
		Get("/files/123456").
		MatchParam("apikey", "api-key").
		Reply(200).
		File("./fixtures/error.json")

	client := threeplay.NewClient("api-key", "secret-key")

	file, err := client.GetFile(123456)
	assert.Equal(err.Error(), "API Error")
	assert.Nil(file)
}

func TestGetFileError(t *testing.T) {
	assert := assert.New(t)

	defer gock.Off()
	gock.New("https://api.3playmedia.com").
		Get("/files/123456").
		MatchParam("apikey", "api-key").
		Reply(200).
		File("./fixtures/not_json")

	client := threeplay.NewClient("api-key", "secret-key")

	file, err := client.GetFile(123456)
	assert.NotNil(err)
	assert.Nil(file)
}

func TestGetFiles(t *testing.T) {
	assert := assert.New(t)

	defer gock.Off()
	gock.New("https://api.3playmedia.com").
		Get("/files").
		MatchParam("apikey", "api-key").
		Reply(200).
		File("./fixtures/files_page1.json")

	client := threeplay.NewClient("api-key", "secret-key")

	filesPage, err := client.GetFiles(nil)
	assert.Nil(err)
	assert.Equal(len(filesPage.Files), 10)
	assert.Equal(filesPage.Summary.CurrentPage.String(), "1")
	assert.Equal(filesPage.Summary.PerPage.String(), "10")
}

func TestFilterFiles(t *testing.T) {
	assert := assert.New(t)
	defer gock.Off()
	gock.New("https://api.3playmedia.com").
		Get("/files").
		MatchParam("apikey", "api-key").
		MatchParam("q", "state=error&video_id=123123").
		Reply(200).
		File("./fixtures/files_page1.json")

	client := threeplay.NewClient("api-key", "secret-key")

	filter := url.Values{
		"video_id": []string{"123123"},
		"state":    []string{"error"},
	}

	filesPage, err := client.FilterFiles(filter, nil)
	assert.Nil(err)
	assert.NotNil(filesPage)
	assert.Equal(len(filesPage.Files), 10)

	filesPage, err = client.FilterFiles(nil, nil)
	assert.Nil(filesPage)
	assert.Equal(err.Error(), "No filters specified")
}

func TestFilterFilesWithPagination(t *testing.T) {
	assert := assert.New(t)
	defer gock.Off()
	gock.New("https://api.3playmedia.com").
		Get("/files").
		MatchParam("apikey", "api-key").
		MatchParam("page", "2").
		MatchParam("per_page", "12").
		MatchParam("q", "state=error&video_id=123123").
		Reply(200).
		File("./fixtures/files_page1.json")

	client := threeplay.NewClient("api-key", "secret-key")

	filter := url.Values{
		"video_id": []string{"123123"},
		"state":    []string{"error"},
	}

	pagination := url.Values{
		"page":     []string{"2"},
		"per_page": []string{"12"},
	}

	filesPage, err := client.FilterFiles(filter, pagination)
	assert.Nil(err)
	assert.NotNil(filesPage)
	assert.Equal(len(filesPage.Files), 10)
}

func TestGetFilesWithPagination(t *testing.T) {
	assert := assert.New(t)

	defer gock.Off()
	gock.New("https://api.3playmedia.com").
		Get("/files").
		MatchParam("apikey", "api-key").
		MatchParam("page", "2").
		Reply(200).
		File("./fixtures/files_page2.json")

	client := threeplay.NewClient("api-key", "secret-key")
	querystring := url.Values{}
	querystring.Add("page", "2")

	filesPage, err := client.GetFiles(querystring)
	assert.Nil(err)
	assert.Equal("2", filesPage.Summary.CurrentPage.String())
}

func TestUploadFileFromURL(t *testing.T) {
	assert := assert.New(t)

	defer gock.Off()
	gock.New("https://api.3playmedia.com").
		Post("/files").
		MatchType("url").
		BodyString("api_secret_key=secret-key&apikey=api-key&link=https%3A%2F%2Fsomewhere.com%2F72397_1_08macron-speech_wg_360p.mp4&video_id=123456").
		Reply(200).
		BodyString("1686514")

	client := threeplay.NewClient("api-key", "secret-key")
	data := url.Values{}
	data.Set("video_id", "123456")

	fileID, err := client.UploadFileFromURL("https://somewhere.com/72397_1_08macron-speech_wg_360p.mp4", data)
	assert.Equal("1686514", fileID)
	assert.Nil(err)
}

func ExampleClient_GetFiles() {
	client := threeplay.NewClient("api-key", "secret")
	filesPage, _ := client.GetFiles(nil)
	fmt.Println(filesPage.Files)

	pagination, _ := url.ParseQuery("page=2&per_page=10")

	filesPage, _ = client.GetFiles(pagination)
	fmt.Println(filesPage.Files)
}

func ExampleClient_UploadFileFromURL() {
	client := threeplay.NewClient("api-key", "secret")
	data, _ := url.ParseQuery("video_id=123&attribute1=abc")
	fileID, _ := client.UploadFileFromURL("http://somewhere.com/video.mp4", data)
	fmt.Println(fileID)
}

func ExampleClient_GetTranscriptWithFormat() {
	client := threeplay.NewClient("api-key", "secret")
	transcript, _ := client.GetTranscriptWithFormat(123, threeplay.JSON)
	fmt.Println(transcript)
}
