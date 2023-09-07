//go:build wireinject
// +build wireinject

package telegram

import (
	"github.com/marlosl/gpt-telegram-bot/utils/config"

	"github.com/google/wire"
)

var TelegramTextSet = wire.NewSet(
	config.ConfigSet,
	newTextService,
)

var TelegramImageSet = wire.NewSet(
	config.ConfigSet,
	newImageService,
)

func NewTextService() (*Telegram, error) {
	wire.Build(TelegramTextSet)
	return &Telegram{}, nil
}

func NewImageService() (*Telegram, error) {
	wire.Build(TelegramImageSet)
	return &Telegram{}, nil
}
