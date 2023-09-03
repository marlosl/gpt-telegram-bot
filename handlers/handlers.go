package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/marlosl/gpt-telegram-bot/clients/sqs"
	"github.com/marlosl/gpt-telegram-bot/consts"
	"github.com/marlosl/gpt-telegram-bot/services/chatgpt"
	"github.com/marlosl/gpt-telegram-bot/services/telegram"
  "github.com/marlosl/gpt-telegram-bot/utils/config"

	"github.com/aws/aws-lambda-go/events"
)

type Chat struct {
	Message string `json:"message"`
}

type Handler struct {
	chatGPT         chatgpt.ChatGPTInterface
	sqsClient       sqs.SQSClientInterface
	telegramService telegram.TelegramInterface
	cfg             config.ConfigInterface
}

func NewHandler() *Handler {
  h := new(Handler)
  h.init()
  return h
}

func (h *Handler) init() {
	if h.chatGPT == nil {
		h.chatGPT = chatgpt.NewChatGPT()
	}

	if h.sqsClient == nil {
		queue := os.Getenv(consts.SendImageQueue)
		h.sqsClient, _ = sqs.NewSQSClient(&queue)
	}

	if h.telegramService == nil {
		h.telegramService = telegram.NewTextService()
	}
}

func (h *Handler) handlePingPong(req events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       "Pong!",
	}, nil
}

func (h *Handler) handleCommandChatTelegram(req events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	var msg telegram.WebhookMessage

	if req.Headers[consts.TelegramWebhookTokenHeader] != h.cfg.TelegramWebhookToken() {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusUnauthorized,
			Body:       "Unauthorized",
		}, nil
	}

	err := h.checkServices()
	if err != nil {
		fmt.Println(err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}, nil
	}

	err = json.Unmarshal([]byte(req.Body), &msg)
	if err != nil {
		fmt.Println(err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}, nil
	}

	updateId := fmt.Sprintf("%d", msg.UpdateId)
	fmt.Printf("Validating if UpdateId exists: %s\n", updateId)

	if h.telegramService.GetCache().ItemExists(updateId) {
		fmt.Println("UpdateId already exists")
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
		}, nil
	}

	fmt.Println("UpdateId does not exist, saving it")
	h.telegramService.GetCache().SaveItem(&updateId)

	command := telegram.GetCommand(&msg.Message.Text)
	fmt.Printf("Command: %s\n", command)

	switch command {
	case telegram.CreateImageCommand:
		return h.handleGenerateImageToTelegram(req, msg, command)
	}
	return h.handleTalkToChatTelegram(req, msg, command)
}

func (h *Handler) handleTalkToChatTelegram(
	req events.APIGatewayV2HTTPRequest,
	msg telegram.WebhookMessage,
	cmd telegram.Command,
) (
	events.APIGatewayProxyResponse,
	error,
) {
	fmt.Println("handleTalkToChatTelegram - start")
	var err error
	var response *chatgpt.ChatResponse

	chatId := ""

	switch cmd {
	case telegram.EditCommand:
		fmt.Println("handleTalkToChatTelegram - EditCommand")
		text, instruction := telegram.ParseMessage(cmd, &msg.Message.Text)
		response, err = h.chatGPT.Edit(*instruction, *text)
	case telegram.None:
		fmt.Println("handleTalkToChatTelegram - None")
		response, err = h.chatGPT.Talk(msg.Message.Text)
	}

	if err != nil {
		fmt.Println(err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}, nil
	}

	if msg.Message != nil && msg.Message.Chat != nil {
		chatId = fmt.Sprintf("%d", msg.Message.Chat.ID)
	}

	if len(response.Choices) == 0 {
		h.telegramService.SendMessage("No Chat GPT response", chatId, false)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusNotFound,
		}, nil
	}

	for _, choice := range response.Choices {
		h.telegramService.SendMessage(choice.Message.Content, chatId, true)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
	}, nil
}

func (h *Handler) handleGenerateImageToTelegram(
	req events.APIGatewayV2HTTPRequest,
	msg telegram.WebhookMessage,
	cmd telegram.Command,
) (
	events.APIGatewayProxyResponse,
	error,
) {
	chatId := ""
	text, _ := telegram.ParseMessage(cmd, &msg.Message.Text)

	response, err := h.chatGPT.CreateImage(*text)
	if err != nil {
		fmt.Println(err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}, nil
	}

	if msg.Message != nil && msg.Message.Chat != nil {
		chatId = fmt.Sprintf("%d", msg.Message.Chat.ID)
	}

	if response == nil || len(response.Data) == 0 {
		h.telegramService.SendMessage("No images were created", chatId, false)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusNotFound,
		}, nil
	}

	h.telegramService.SendMessage(fmt.Sprintf("Number of generated images: %d", len(response.Data)), chatId, true)

	for i, photo := range response.Data {
		if h.cfg.SendImageByUrl() {
			h.telegramService.SendMessage(fmt.Sprintf("Image url: %s", photo.Url), chatId, true)
		} else {
			h.sendPhoto(h.telegramService, i, photo.Url, chatId)
		}
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
	}, nil
}

func (h *Handler) sendPhoto(t telegram.TelegramInterface, i int, url string, chatId string) {
	fmt.Printf("Sending image %d: %s\n", i, url)
	message := &telegram.ImageMessage{
		ChatId:   chatId,
		ImageUrl: url,
	}

	err := h.sqsClient.SendMsg(message)
	if err != nil {
		t.SendMessage(fmt.Sprintf("Error while sending image: %v\n", err), chatId, true)
	}
}

func (h *Handler) handleTalkToChatGPT(req events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	var chat Chat
	err := json.Unmarshal([]byte(req.Body), &chat)
	if err != nil {
		fmt.Println(err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}, nil
	}

	response, err := h.chatGPT.Talk(chat.Message)
	if err != nil {
		fmt.Println(err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}, nil
	}

	if len(response.Choices) == 0 {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusNotFound,
			Body:       "No Chat GPT response",
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       response.Choices[0].Message.Content,
	}, nil
}

func (h *Handler) checkServices() error {
	if h.chatGPT == nil {
		return errors.New("chatGPT is not initialized")
	}

	if h.sqsClient == nil {
		return errors.New("sqsClient is not initialized")
	}

	if h.telegramService == nil {
		return errors.New("telegramService is not initialized")
	}

	return nil
}
