package handlers

import (
	"fmt"
	"net/http"

	"github.com/marlosl/gpt-telegram-bot/utils"

	"github.com/aws/aws-lambda-go/events"
)

func Router(req events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Printf("Request: %s", utils.SPrintJson(req))
	if req.RequestContext.HTTP.Method == "GET" {

		if req.RawPath == "/ping" {
			return handlePingPong(req)
		}
	}
	if req.RequestContext.HTTP.Method == "POST" {
		if req.RawPath == "/gpt" {
			return handleTalkToChatGPT(req)
		}
		if req.RawPath == "/telegram-bot" {
			return handleCommandChatTelegram(req)
		}
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusMethodNotAllowed,
		Body:       http.StatusText(http.StatusMethodNotAllowed),
	}, nil
}
