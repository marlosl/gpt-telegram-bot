package chatgpt

import (
	"fmt"
	"strings"
	"time"

	"github.com/marlosl/gpt-telegram-bot/utils"
	"github.com/marlosl/gpt-telegram-bot/utils/config"

	"github.com/go-resty/resty/v2"
)

type ChatGPT struct {
	ApiKey         string
	GptModel       string
	ChatUrl        string
	CreateImageUrl string
	EditUrl        string
}

func (c *ChatGPT) InitApi() {
	config.NewConfig(config.SSM)
	c.ApiKey = config.Store.GptApiKey
	c.GptModel = config.Store.GptModel
	c.ChatUrl = "https://api.openai.com/v1/chat/completions"
	c.CreateImageUrl = "https://api.openai.com/v1/images/generations"
	c.EditUrl = "https://api.openai.com/v1/edits"
}

func (c *ChatGPT) CreateRequest() *resty.Request {
	client := resty.New()
	client.SetTimeout(5 * time.Minute)
	return client.R().
		SetAuthToken(c.ApiKey).
		SetHeader("accept", "*/*").
		SetHeader("accept-encoding", "gzip, deflate, br").
		SetHeader("accept-language", "pt-BR,pt;q=0.9,en-US;q=0.8,en;q=0.7").
		SetHeader("cache-control", "no-cache").
		SetHeader("content-type", "application/json").
		SetHeader("pragma", "no-cache").
		SetHeader("sec-ch-ua", "\"Google Chrome\";v=\"105\", \"Not)A;Brand\";v=\"8\", \"Chromium\";v=\"105\"").
		SetHeader("sec-ch-ua-platform", "\"macOS\"").
		SetHeader("sec-fetch-dest", "empty").
		SetHeader("sec-fetch-mode", "cors").
		SetHeader("sec-fetch-site", "same-site").
		SetHeader("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36").
		EnableTrace()
}

func (c *ChatGPT) isSuccess(r *resty.Response) bool {
	return r != nil && r.StatusCode() >= 200 && r.StatusCode() <= 299
}

func (c *ChatGPT) Talk(message string) (*ChatResponse, error) {
	fmt.Printf("Talk: %s\n", message)
	resp, err := c.CreateRequest().
		SetResult(ChatResponse{}).
		SetBody(c.CreateChatRequest(message)).
		Post(c.ChatUrl)

	utils.PrintRestyDebug(resp, err)
	if err != nil || !c.isSuccess(resp) {

		return nil, err
	}
	return resp.Result().(*ChatResponse), nil
}

func (c *ChatGPT) Edit(instruction, message string) (*ChatResponse, error) {
	resp, err := c.CreateRequest().
		SetResult(ChatResponse{}).
		SetBody(c.CreateEditRequest(instruction, message)).
		Post(c.ChatUrl)

	utils.PrintRestyDebug(resp, err)
	if err != nil || !c.isSuccess(resp) {

		return nil, err
	}
	return resp.Result().(*ChatResponse), nil
}

func (c *ChatGPT) CreateChatRequest(message string) ChatRequest {
	var stop string
	var req ChatRequest

	if strings.Contains(message, "\"\"\"") {
		stop = "\"\"\""
	}

	if strings.Contains(message, "###") {
		stop = "###"
	}

	if len(stop) > 0 {
		req = ChatRequest{
			Model: c.GptModel,
			Messages: []ChatMessage{
				{
					Role:    "user",
					Content: message,
				},
			},
			Stop: &stop,
		}
	} else {
		req = ChatRequest{
			Model: c.GptModel,
			Messages: []ChatMessage{
				{
					Role:    "user",
					Content: message,
				},
			},
		}
	}

	fmt.Printf("ChatRequest: %s\n", utils.SPrintJson(req))
	return req
}

func (c *ChatGPT) CreateEditRequest(instruction, message string) EditRequest {
	return EditRequest{
		Model:       "text-davinci-edit-001",
		Input:       message,
		Instruction: instruction,
	}
}

func (c *ChatGPT) CreateImage(message string) (*CreateImageResponse, error) {
	resp, err := c.CreateRequest().
		SetResult(CreateImageResponse{}).
		SetBody(c.CreateImageRequest(message)).
		Post(c.CreateImageUrl)

	utils.PrintRestyDebug(resp, err)
	if err != nil || !c.isSuccess(resp) {

		return nil, err
	}
	return resp.Result().(*CreateImageResponse), nil
}

func (c *ChatGPT) CreateImageRequest(message string) CreateImageRequest {
	return CreateImageRequest{
		Prompt: message,
		N:      2,
		Size:   "1024x1024",
	}
}

func NewChatGPT() *ChatGPT {
	c := &ChatGPT{}
	c.InitApi()
	return c
}
