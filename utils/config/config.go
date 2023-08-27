package config

import (
	"os"
	"sync"

	"github.com/marlosl/gpt-telegram-bot/clients/ssm"
	"github.com/marlosl/gpt-telegram-bot/consts"
)

type ConfigType int64

const (
	SSM  ConfigType = 0
	File ConfigType = 1
)

var (
	mutex = &sync.Mutex{}
	Store *Config
)

type Config struct {
	TelegramBotTextToken  string
	TelegramBotImageToken string
	GptApiKey             string
	SendImageByUrl        bool
	GptModel              string
	TelegramWebhookToken  string
}

func NewConfig(t ConfigType) *Config {
	if Store == nil {
		mutex.Lock()
		defer mutex.Unlock()
		if Store == nil {
			switch t {
			case SSM:
				Store = &Config{
					TelegramBotTextToken:  ssm.Get(consts.PARAMETER_TELEGRAM_BOT_TEXT_TOKEN),
					TelegramBotImageToken: ssm.Get(consts.PARAMETER_TELEGRAM_BOT_IMAGE_TOKEN),
					GptApiKey:             ssm.Get(consts.PARAMETER_GPT_API_KEY),
					SendImageByUrl:        ssm.Get(consts.PARAMETER_SEND_IMAGE_BY_URL) == "true",
					GptModel:              ssm.Get(consts.PARAMETER_GPT_MODEL),
					TelegramWebhookToken:  ssm.Get(consts.PARAMETER_TELEGRAM_WEBHOOK_TOKEN),
				}
			case File:
				Store = &Config{
					TelegramBotTextToken:  os.Getenv(consts.TelegramBotTextToken),
					TelegramBotImageToken: os.Getenv(consts.TelegramBotImageToken),
					GptApiKey:             os.Getenv(consts.GptApiKey),
					SendImageByUrl:        false,
					GptModel:              "",
					TelegramWebhookToken:  "",
				}
			}
		}
	}
	return Store
}

func NewSSMConfig() *Config {
  return NewConfig(SSM)
}

func NewFileConfig() *Config {
  return NewConfig(File)
}
