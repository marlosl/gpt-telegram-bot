package main

import (
	"github.com/marlosl/gpt-telegram-bot/handlers"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handlers.Router)
}
