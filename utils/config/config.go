package config

import (
	"os"

	"github.com/marlosl/gpt-telegram-bot/clients/ssm"
	"github.com/marlosl/gpt-telegram-bot/consts"
)

type ConfigType int64

const (
	SSM  ConfigType = 0
	File ConfigType = 1
)

type ConfigInterface interface {
	GetTelegramBotTextToken() string
	GetTelegramBotImageToken() string
	GetGptApiKey() string
	GetSendImageByUrl() bool
	GetGptModel() string
	GetTelegramWebhookToken() string
}
type Config struct {
	telegramBotTextToken  string
	telegramBotImageToken string
	telegramWebhookToken  string
	sendImageByUrl        bool
	gptApiKey             string
	gptModel              string

	ssm ssm.SSMClientInterface
}

func NewConfig(t ConfigType, ssm *ssm.SSMClientInterface) *Config {
	c := new(Config)

	if t == SSM {
		c.ssm = *ssm
	}

	c.init(t)
	return c
}

func NewSSMConfig(ssm *ssm.SSMClientInterface) *Config {
	return NewConfig(SSM, ssm)
}

func NewFileConfig() *Config {
	return NewConfig(File, nil)
}

func (c *Config) init(t ConfigType) {
	switch t {
	case SSM:
		c.telegramBotTextToken = c.ssm.Get(consts.PARAMETER_TELEGRAM_BOT_TEXT_TOKEN)
		c.telegramBotImageToken = c.ssm.Get(consts.PARAMETER_TELEGRAM_BOT_IMAGE_TOKEN)
		c.gptApiKey = c.ssm.Get(consts.PARAMETER_GPT_API_KEY)
		c.sendImageByUrl = c.ssm.Get(consts.PARAMETER_SEND_IMAGE_BY_URL) == "true"
		c.gptModel = c.ssm.Get(consts.PARAMETER_GPT_MODEL)
		c.telegramWebhookToken = c.ssm.Get(consts.PARAMETER_TELEGRAM_WEBHOOK_TOKEN)
	case File:
		c.telegramBotTextToken = os.Getenv(consts.TelegramBotTextToken)
		c.telegramBotImageToken = os.Getenv(consts.TelegramBotImageToken)
		c.gptApiKey = os.Getenv(consts.GptApiKey)
		c.sendImageByUrl = false
		c.gptModel = ""
		c.telegramWebhookToken = ""
	}
}

func (c *Config) GetTelegramBotTextToken() string {
	return c.telegramBotTextToken
}

func (c *Config) GetTelegramBotImageToken() string {
	return c.telegramBotImageToken
}

func (c *Config) GetGptApiKey() string {
	return c.gptApiKey
}

func (c *Config) GetSendImageByUrl() bool {
	return c.sendImageByUrl
}

func (c *Config) GetGptModel() string {
	return c.gptModel
}

func (c *Config) GetTelegramWebhookToken() string {
	return c.telegramWebhookToken
}
