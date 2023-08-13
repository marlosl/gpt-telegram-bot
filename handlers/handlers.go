package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/marlosl/gpt-telegram-bot/config"
	"github.com/marlosl/gpt-telegram-bot/consts"
	"github.com/marlosl/gpt-telegram-bot/services/chatgpt"
	"github.com/marlosl/gpt-telegram-bot/services/sqs"
	"github.com/marlosl/gpt-telegram-bot/services/telegram"

	"github.com/aws/aws-lambda-go/events"
)

type Chat struct {
	Message string `json:"message"`
}

var (
	chatGPT         *chatgpt.ChatGPT
	sqsClient       *sqs.SQSClient
	telegramService *telegram.Telegram
)

func init() {
	config.NewConfig(config.SSM)
	if chatGPT == nil {
		chatGPT = chatgpt.NewChatGPT()
	}

	if sqsClient == nil {
		queue := os.Getenv(consts.SendImageQueue)
		sqsClient, _ = sqs.NewSQSClient(&queue)
	}

	if telegramService == nil {
		telegramService = telegram.NewTextService()
	}
}

func handlePingPong(req events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       "Pong!",
	}, nil
}

func handleCommandChatTelegram(req events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	var msg telegram.WebhookMessage

	if req.Headers[consts.TelegramWebhookTokenHeader] != config.Store.TelegramWebhookToken {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusUnauthorized,
			Body:       "Unauthorized",
		}, nil
	}

	err := checkServices()
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

	if telegramService.Cache.ItemExists(updateId) {
		fmt.Println("UpdateId already exists")
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
		}, nil
	}

	fmt.Println("UpdateId does not exist, saving it")
	telegramService.Cache.SaveItem(&updateId)

	command := telegram.GetCommand(&msg.Message.Text)
	fmt.Printf("Command: %s\n", command)

	switch command {
	case telegram.CreateImageCommand:
		return handleGenerateImageToTelegram(req, msg, command)
	}
	return handleTalkToChatTelegram(req, msg, command)
}

func handleTalkToChatTelegram(
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
		response, err = chatGPT.Edit(*instruction, *text)
	case telegram.None:
		fmt.Println("handleTalkToChatTelegram - None")
		response, err = chatGPT.Talk(msg.Message.Text)
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
		telegramService.SendMessage("No Chat GPT response", chatId, false)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusNotFound,
		}, nil
	}

	for _, choice := range response.Choices {
		telegramService.SendMessage(choice.Message.Content, chatId, true)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
	}, nil
}

func handleGenerateImageToTelegram(
	req events.APIGatewayV2HTTPRequest,
	msg telegram.WebhookMessage,
	cmd telegram.Command,
) (
	events.APIGatewayProxyResponse,
	error,
) {
	chatId := ""
	text, _ := telegram.ParseMessage(cmd, &msg.Message.Text)

	response, err := chatGPT.CreateImage(*text)
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
		telegramService.SendMessage("No images were created", chatId, false)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusNotFound,
		}, nil
	}

	telegramService.SendMessage(fmt.Sprintf("Number of generated images: %d", len(response.Data)), chatId, true)

	for i, photo := range response.Data {
		if config.Store.SendImageByUrl {
			telegramService.SendMessage(fmt.Sprintf("Image url: %s", photo.Url), chatId, true)
		} else {
			sendPhoto(telegramService, i, photo.Url, chatId)
		}
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
	}, nil
}

func sendPhoto(t *telegram.Telegram, i int, url string, chatId string) {
	fmt.Printf("Sending image %d: %s\n", i, url)
	message := &telegram.ImageMessage{
		ChatId:   chatId,
		ImageUrl: url,
	}

	err := sqsClient.SendMsg(message)
	if err != nil {
		t.SendMessage(fmt.Sprintf("Error while sending image: %v\n", err), chatId, true)
	}
}

func handleTalkToChatGPT(req events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	var chat Chat
	err := json.Unmarshal([]byte(req.Body), &chat)
	if err != nil {
		fmt.Println(err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}, nil
	}

	response, err := chatGPT.Talk(chat.Message)
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

func checkServices() error {
	if chatGPT == nil {
		return errors.New("chatGPT is not initialized")
	}

	if sqsClient == nil {
		return errors.New("sqsClient is not initialized")
	}

	if telegramService == nil {
		return errors.New("telegramService is not initialized")
	}

	return nil
}
