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

func NewTextService() *Telegram {
	wire.Build(TelegramTextSet)
	return Telegram{}
}

func NewImageService() *Telegram {
	wire.Build(TelegramImageSet)
	return Telegram{}
}
