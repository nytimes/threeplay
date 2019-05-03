package threeplay_test

import (
	"testing"

	"github.com/nytimes/threeplay"
	"github.com/stretchr/testify/assert"
	gock "gopkg.in/h2non/gock.v1"
)

func TestGetTags(t *testing.T) {
	assert := assert.New(t)

	defer gock.Off()
	gock.New("https://api.3playmedia.com").
		Get("/files/123456").
		MatchParam("apikey", "api-key").
		Reply(200).
		BodyString(`["physics","robots","spycraft"]`)

	client := threeplay.NewClient("api-key", "secret-key")
	tags, err := client.GetTags(123456)
	assert.Equal("physics", tags[0])
	assert.Nil(err)
}

func TestGetTagsAPIError(t *testing.T) {
	assert := assert.New(t)

	defer gock.Off()
	gock.New("https://api.3playmedia.com").
		Get("/files/123456").
		MatchParam("apikey", "api-key").
		Reply(200).
		File("./fixtures/error.json")

	client := threeplay.NewClient("api-key", "secret-key")
	tags, err := client.GetTags(123456)
	assert.Nil(tags)
	assert.NotNil(err)
	assert.Equal(threeplay.ErrUnauthorized.Error(), err.Error())
}

func TestGetTagsError(t *testing.T) {
	assert := assert.New(t)

	defer gock.Off()
	gock.New("https://api.3playmedia.com").
		Get("/files/123456").
		MatchParam("apikey", "api-key").
		Reply(200).
		File("./fixtures/not_json")

	client := threeplay.NewClient("api-key", "secret-key")

	tags, err := client.GetTags(123456)
	assert.NotNil(err)
	assert.Nil(tags)
}

func TestAddTag(t *testing.T) {
	assert := assert.New(t)

	defer gock.Off()
	gock.New("https://api.3playmedia.com").
		Post("/files/123456/tags").
		BodyString("api_secret_key=secret-key&apikey=api-key&name=this-is-a-tag").
		Reply(200).
		File("./fixtures/add_tag.json")

	client := threeplay.NewClient("api-key", "secret-key")

	tags, err := client.AddTag(123456, "this-is-a-tag")
	assert.Nil(err)
	assert.NotNil(tags)
	assert.Equal("this-is-a-tag", tags[0])
}

func TestAddTagApiError(t *testing.T) {
	assert := assert.New(t)

	defer gock.Off()
	gock.New("https://api.3playmedia.com").
		Post("/files/123456/tags").
		BodyString("api_secret_key=secret-key&apikey=api-key&name=this-is-a-tag").
		Reply(200).
		File("./fixtures/error.json")

	client := threeplay.NewClient("api-key", "secret-key")

	tags, err := client.AddTag(123456, "this-is-a-tag")
	assert.Nil(tags)
	assert.NotNil(err)
	assert.Equal(threeplay.ErrUnauthorized.Error(), err.Error())
}

func TestAddTagError(t *testing.T) {
	assert := assert.New(t)

	defer gock.Off()
	gock.New("https://api.3playmedia.com").
		Post("/files/123456/tags").
		BodyString("api_secret_key=secret-key&apikey=api-key&name=this-is-a-tag").
		Reply(200).
		File("./fixtures/not_json")

	client := threeplay.NewClient("api-key", "secret-key")

	tags, err := client.AddTag(123456, "this-is-a-tag")
	assert.Nil(tags)
	assert.NotNil(err)
}

func TestRemoveTag(t *testing.T) {
	assert := assert.New(t)

	defer gock.Off()
	gock.New("https://api.3playmedia.com").
		Post("/files/123456/tags/this-is-a-tag").
		MatchType("url").
		BodyString("_method=delete&api_secret_key=secret-key&apikey=api-key").
		Reply(200).
		BodyString(`["physics","robots","spycraft"]`)

	client := threeplay.NewClient("api-key", "secret-key")

	tags, err := client.RemoveTag(123456, "this-is-a-tag")
	assert.Nil(err)
	assert.NotNil(tags)
	assert.Equal("physics", tags[0])
}

func TestRemoveTagApiError(t *testing.T) {
	assert := assert.New(t)

	defer gock.Off()
	gock.New("https://api.3playmedia.com").
		Post("/files/123456/tags/this-is-a-tag").
		MatchType("url").
		BodyString("_method=delete&api_secret_key=secret-key&apikey=api-key").
		Reply(200).
		File("./fixtures/error.json")

	client := threeplay.NewClient("api-key", "secret-key")

	tags, err := client.RemoveTag(123456, "this-is-a-tag")
	assert.Nil(tags)
	assert.NotNil(err)
	assert.Equal(threeplay.ErrUnauthorized.Error(), err.Error())
}

func TestRemoveTagError(t *testing.T) {
	assert := assert.New(t)

	defer gock.Off()
	gock.New("https://api.3playmedia.com").
		Post("/files/123456/tags/this-is-a-tag").
		MatchType("url").
		BodyString("_method=delete&api_secret_key=secret-key&apikey=api-key").
		Reply(200).
		File("./fixtures/not_json")

	client := threeplay.NewClient("api-key", "secret-key")

	tags, err := client.RemoveTag(123456, "this-is-a-tag")
	assert.Nil(tags)
	assert.NotNil(err)
}
