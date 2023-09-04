//go:build wireinject
// +build wireinject

package chatgpt

import (
	"github.com/marlosl/gpt-telegram-bot/utils/config"

	"github.com/google/wire"
)

var ChatGPTSet = wire.NewSet(
	config.ConfigSet,
	newChatGPT,
	wire.Bind(new(ChatGPTInterface), new(*ChatGPT)),
)

func NewChatGPT() *ChatGPT {
	wire.Build(ChatGPTSet)
	return &ChatGPT{}
}
