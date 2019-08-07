package v3api_test

import (
	"testing"

	"github.com/nytimes/threeplay/v3api"
	"github.com/stretchr/testify/assert"
	gock "gopkg.in/h2non/gock.v1"
)

func TestOrderTranscript(t *testing.T) {
	assert := assert.New(t)
	defer gock.Off()

	gock.New("https://api.3playmedia.com").
		Post("/v3/transcripts/order/asr").
		MatchType("url").
		BodyString("api_key=api-key&media_file_id=3628518").
		Reply(200).
		File("../fixtures/v3_transcript_order_200.json")

	client := v3api.NewClient("api-key")

	callParams := v3api.CallParams{APIKey: ""}
	transcriptData, err := client.OrderTranscript("3628518", "", "asr", callParams)
	assert.Nil(err)
	assert.NotNil(transcriptData)
	assert.Equal("pending", transcriptData.Status)
	assert.Equal("AsrTranscript", transcriptData.Type)
}

func TestOrderTranscriptCustomAPIKey(t *testing.T) {
	assert := assert.New(t)
	defer gock.Off()

	gock.New("https://api.3playmedia.com").
		Post("/v3/transcripts/order/asr").
		MatchType("url").
		BodyString("api_key=custom-key&media_file_id=3628518").
		Reply(200).
		File("../fixtures/v3_transcript_order_200.json")

	client := v3api.NewClient("api-key")

	callParams := v3api.CallParams{APIKey: "custom-key"}
	transcriptData, err := client.OrderTranscript("3628518", "", "asr", callParams)
	assert.Nil(err)
	assert.NotNil(transcriptData)
	assert.Equal("pending", transcriptData.Status)
	assert.Equal("AsrTranscript", transcriptData.Type)
}

func TestOrderTranscriptError(t *testing.T) {
	assert := assert.New(t)
	defer gock.Off()

	gock.New("https://api.3playmedia.com").
		Post("/v3/transcripts/order/asr").
		MatchType("url").
		BodyString("api_key=api-key&media_file_id=123456").
		Reply(404).
		File("../fixtures/v3_transcript_order_404.json")

	client := v3api.NewClient("api-key")

	callParams := v3api.CallParams{APIKey: ""}
	transcriptData, err := client.OrderTranscript("123456", "", "asr", callParams)
	assert.Empty(transcriptData)
	assert.NotNil(err)
	assert.Equal("404: not_found_error-Not found", err.Error())
}

func TestTranscriptInfo(t *testing.T) {
	assert := assert.New(t)
	defer gock.Off()

	gock.New("https://api.3playmedia.com").
		Get("/v3/transcripts/3633088").
		MatchParam("api_key", "api-key").
		Reply(200).
		File("../fixtures/v3_transcript_info_complete.json")

	client := v3api.NewClient("api-key")

	callParams := v3api.CallParams{APIKey: ""}
	transcriptData, err := client.GetTranscriptInfo("3633088", callParams)
	assert.Nil(err)
	assert.NotNil(transcriptData)
	assert.Equal("complete", transcriptData.Status)
}

func TestTranscriptInfoError(t *testing.T) {
	assert := assert.New(t)
	defer gock.Off()

	gock.New("https://api.3playmedia.com").
		Get("/v3/transcripts/123").
		MatchParam("api_key", "api-key").
		Persist().
		Reply(500).
		File("../fixtures/v3_unknown_error.json")

	client := v3api.NewClient("api-key")

	callParams := v3api.CallParams{APIKey: ""}
	transcriptData, err := client.GetTranscriptInfo("123", callParams)
	assert.Empty(transcriptData)
	assert.NotNil(err)
	assert.Equal("500: standard_error-Internal server error", err.Error())
}

func TestTranscriptText(t *testing.T) {
	assert := assert.New(t)

	defer gock.Off()
	gock.New("https://api.3playmedia.com").
		Get("/v3/transcripts/3633088/text").
		MatchParam("api_key", "api-key").
		MatchParam("output_format_id", "139").
		Reply(200).
		File("../fixtures/v3_transcript_text.json")

	client := v3api.NewClient("api-key")

	callParams := v3api.CallParams{APIKey: ""}
	transcript, err := client.GetTranscriptText("3633088", "", "vtt", callParams)
	assert.Nil(err)
	assert.NotEmpty(transcript)
}

func TestTranscriptTextError(t *testing.T) {
	assert := assert.New(t)

	defer gock.Off()
	gock.New("https://api.3playmedia.com").
		Get("/v3/transcripts/9846/text").
		MatchParam("api_key", "api-key").
		MatchParam("output_format_id", "139").
		Persist().
		Reply(500).
		File("../fixtures/v3_unknown_error.json")

	client := v3api.NewClient("api-key")

	callParams := v3api.CallParams{APIKey: ""}
	transcript, err := client.GetTranscriptText("9846", "", "vtt", callParams)
	assert.Empty(transcript)
	assert.NotNil(err)
	assert.Equal("500: standard_error-Internal server error", err.Error())
}

func TestTranscriptCancel(t *testing.T) {
	assert := assert.New(t)
	defer gock.Off()

	gock.New("https://api.3playmedia.com").
		Post("/transcripts/8794567/cancel").
		MatchType("url").
		BodyString("api_key=api-key").
		Reply(200).
		File("../fixtures/v3_cancel_200.json")

	client := v3api.NewClient("api-key")

	callParams := v3api.CallParams{APIKey: ""}
	err := client.CancelTranscript("8794567", callParams)
	assert.Nil(err)
}

func TestTranscriptCancelError(t *testing.T) {
	assert := assert.New(t)
	defer gock.Off()

	gock.New("https://api.3playmedia.com").
		Post("/transcripts/843759/cancel").
		MatchType("url").
		BodyString("api_key=api-key").
		Reply(403).
		File("../fixtures/v3_cancel_403.json")

	client := v3api.NewClient("api-key")

	callParams := v3api.CallParams{APIKey: ""}
	err := client.CancelTranscript("843759", callParams)
	assert.NotNil(err)
	assert.Equal("403: forbidden_error-You cannot cancel this transcript at this time.", err.Error())
}

func TestTranscriptEditingLink(t *testing.T) {
	assert := assert.New(t)

	defer gock.Off()
	gock.New("https://api.3playmedia.com").
		Get("/v3/transcripts/3633088/expiring_editing_link").
		MatchParam("api_key", "api-key").
		MatchParam("hours_until_expiration", "2").
		Reply(200).
		File("../fixtures/v3_transcript_link_200.json")

	client := v3api.NewClient("api-key")

	callParams := v3api.CallParams{APIKey: ""}
	link, err := client.GetEditingLink("3633088", 2, callParams)
	assert.Nil(err)
	assert.Equal("http://external.3playmedia.com/transcripts/3633088/edit?exp_key=key", link)
}

func TestTranscriptEditingLinkError(t *testing.T) {
	assert := assert.New(t)

	defer gock.Off()
	gock.New("https://api.3playmedia.com").
		Get("/v3/transcripts/bad-id/expiring_editing_link").
		MatchParam("api_key", "api-key").
		MatchParam("hours_until_expiration", "2").
		Reply(404).
		File("../fixtures/v3_transcript_order_404.json")

	client := v3api.NewClient("api-key")

	callParams := v3api.CallParams{APIKey: ""}
	link, err := client.GetEditingLink("bad-id", 2, callParams)
	assert.Empty(link)
	assert.NotNil(err)
	assert.Equal("404: not_found_error-Not found", err.Error())
}
