package chatgpt

import (
	"net/http"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

func TestInitApi(t *testing.T) {
	chat := NewChatGPT()

	assert.NotNil(t, chat.ApiKey)
	assert.Equal(t, "https://api.openai.com/v1/chat/completions", chat.ChatUrl)
}

func TestCreateRequest(t *testing.T) {
	chat := NewChatGPT()
	request := chat.CreateRequest()

	assert.NotNil(t, request)
	assert.Equal(t, chat.ApiKey, request.Header.Get("Authorization")[7:]) // Removing "Bearer " prefix
}

func TestIsSuccess(t *testing.T) {
	chat := NewChatGPT()

	assert.True(t, chat.isSuccess(&resty.Response{StatusCode: http.StatusOK}))
	assert.False(t, chat.isSuccess(&resty.Response{StatusCode: http.StatusBadRequest}))
}

func TestTalk(t *testing.T) {
	chat := NewChatGPT()

	// Mock the response
	resty.DefaultClient.SetMock(&resty.MockResponse{
		Error:    nil,
		Response: http.Response{StatusCode: http.StatusOK},
		Body:     `{"id": "abcd", "content": "Hello!"}`,
	})

	resp, err := chat.Talk("Hello ChatGPT!")
	assert.Nil(t, err)
	assert.Equal(t, "Hello!", resp.Content)
}

func TestEdit(t *testing.T) {
	chat := NewChatGPT()

	// Mock the response
	resty.DefaultClient.SetMock(&resty.MockResponse{
		Error:    nil,
		Response: http.Response{StatusCode: http.StatusOK},
		Body:     `{"id": "editabcd", "content": "Edited content"}`,
	})

	resp, err := chat.Edit("Edit this", "Original content")
	assert.Nil(t, err)
	assert.Equal(t, "Edited content", resp.Content)
}

func TestCreateImage(t *testing.T) {
	chat := NewChatGPT()

	// Mock the response
	resty.DefaultClient.SetMock(&resty.MockResponse{
		Error:    nil,
		Response: http.Response{StatusCode: http.StatusOK},
		Body:     `{"id": "imageabcd", "url": "https://fakeimageurl.com/image.jpg"}`,
	})

	resp, err := chat.CreateImage("Create an image of a sunset.")
	assert.Nil(t, err)
	assert.Equal(t, "https://fakeimageurl.com/image.jpg", resp.URL)
}

// ... Additional tests for other utility functions and edge cases.
