package threeplay_test

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/NYTimes/threeplay"
	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
)

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
	assert.Equal(err.Error(), "Must specify new data")

	gock.New("https://api.3playmedia.com").
		Put("/files/123456").
		MatchType("url").
		BodyString("api_secret_key=secret-key&apikey=api-key&name=other-name").
		Reply(200).
		File("./fixtures/error.json")
	data, _ := url.ParseQuery("name=other-name")

	err = client.UpdateFile(123456, data)

	assert.NotNil(err)
	assert.Equal(err.Error(), "Api Error")
}

func ExampleUpdateFile() {
	client := threeplay.NewClient("api-key", "api-secret")
	data, _ := url.ParseQuery("name=another-name")
	err := client.UpdateFile(1687446, data)
	if err != nil {
		fmt.Println(err)
	}
}
