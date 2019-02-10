package threeplay_test

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/NYTimes/threeplay"
	"github.com/stretchr/testify/assert"
	gock "gopkg.in/h2non/gock.v1"
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
	assert.Equal(threeplay.ErrUnauthorized.Error(), err.Error())
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

	filesPage, err := client.GetFiles(nil, nil)
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

	filesPage, err := client.GetFiles(nil, filter)
	assert.Nil(err)
	assert.NotNil(filesPage)
	assert.Equal(len(filesPage.Files), 10)
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

	filesPage, err := client.GetFiles(pagination, filter)
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

	filesPage, err := client.GetFiles(querystring, nil)
	assert.Nil(err)
	assert.Equal("2", filesPage.Summary.CurrentPage.String())
}

func TestUpdateFile(t *testing.T) {
	assert := assert.New(t)
	defer gock.Off()

	gock.New("https://api.3playmedia.com").
		Put("/files/123456").
		MatchType("url").
		BodyString("api_secret_key=secret-key&apikey=api-key&name=other-name").
		Reply(200).
		BodyString("1")

	client := threeplay.NewClient("api-key", "secret-key")
	data, _ := url.ParseQuery("name=other-name")
	err := client.UpdateFile(123456, data)
	assert.Nil(err)
}

func TestUpdateFileError(t *testing.T) {
	assert := assert.New(t)
	defer gock.Off()

	client := threeplay.NewClient("api-key", "secret-key")
	err := client.UpdateFile(123456, nil)
	assert.NotNil(err)
	assert.Equal(err.Error(), "must specify new data")

	gock.New("https://api.3playmedia.com").
		Put("/files/123456").
		MatchType("url").
		BodyString("api_secret_key=secret-key&apikey=api-key&name=other-name").
		Reply(200).
		File("./fixtures/error.json")
	data, _ := url.ParseQuery("name=other-name")

	err = client.UpdateFile(123456, data)

	assert.NotNil(err)
	assert.Equal(threeplay.ErrUnauthorized.Error(), err.Error())

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
	assert.Equal(uint(1686514), fileID)
	assert.Nil(err)
}

func TestUploadFileFromURLInvalidResponse(t *testing.T) {
	assert := assert.New(t)

	defer gock.Off()
	gock.New("https://api.3playmedia.com").
		Post("/files").
		MatchType("url").
		BodyString("api_secret_key=secret-key&apikey=api-key&link=https%3A%2F%2Fsomewhere.com%2F72397_1_08macron-speech_wg_360p.mp4&video_id=123456").
		Reply(200).
		BodyString("<p>Something went wrong, but I still return 200!</p>")

	client := threeplay.NewClient("api-key", "secret-key")
	data := url.Values{}
	data.Set("video_id", "123456")

	fileID, err := client.UploadFileFromURL("https://somewhere.com/72397_1_08macron-speech_wg_360p.mp4", data)
	assert.Equal(uint(0), fileID)
	assert.NotNil(err)
	assert.Equal(err.Error(), "invalid response: <p>Something went wrong, but I still return 200!</p>")
}

func ExampleClient_GetFiles() {
	client := threeplay.NewClient("api-key", "secret")
	filesPage, _ := client.GetFiles(nil, nil)
	fmt.Println(filesPage.Files)

	pagination, _ := url.ParseQuery("page=2&per_page=10")

	filesPage, _ = client.GetFiles(pagination, nil)
	fmt.Println(filesPage.Files)
}

func ExampleClient_UploadFileFromURL() {
	client := threeplay.NewClient("api-key", "secret")
	data, _ := url.ParseQuery("video_id=123&attribute1=abc")
	fileID, _ := client.UploadFileFromURL("http://somewhere.com/video.mp4", data)
	fmt.Println(fileID)
}

func ExampleClient_UpdateFile() {
	client := threeplay.NewClient("api-key", "api-secret")
	data, _ := url.ParseQuery("name=another-name")
	err := client.UpdateFile(1687446, data)
	if err != nil {
		fmt.Println(err)
	}
}
