package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/marlosl/gpt-telegram-bot//utils/config"
	"github.com/marlosl/gpt-telegram-bot/services/telegram"
	"github.com/marlosl/gpt-telegram-bot/utils"

	"github.com/aws/aws-lambda-go/events"
)

type SQSMessage struct {
	Message string `json:"Message"`
}

func SendImageHandler(ctx context.Context, sqsEvent events.SQSEvent) error {
	config.NewConfig(config.SSM)
	telegramService := telegram.NewTextService()
	for _, message := range sqsEvent.Records {
		fmt.Printf("The message %s for event source %s = %s \n", message.MessageId, message.EventSource, message.Body)

		var imgMsg telegram.ImageMessage
		err := json.Unmarshal([]byte(message.Body), &imgMsg)
		if err != nil {
			fmt.Printf("Can't unmarshal sqsMessage: %v \n", err)
			continue
		}

		fmt.Printf("Unmarshal signal: %s \n", utils.SPrintJson(imgMsg))

		if imgMsg.ImageUrl == "" {
			fmt.Println("ImageUrl is empty")
			continue
		}

		err = telegramService.SendPhotoGet(imgMsg.ImageUrl, imgMsg.ChatId)
		if err != nil {
			telegramService.SendMessage(fmt.Sprintf("Error while sending image: %v\n", err), imgMsg.ChatId, true)
		}
		fmt.Println("Image sent")
	}

	return nil
}
