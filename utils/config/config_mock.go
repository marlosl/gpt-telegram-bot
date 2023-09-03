package config

import (
	"github.com/stretchr/testify/mock"
)

type ConfigMock struct {
	mock.Mock
}

func (c *ConfigMock) GetTelegramBotTextToken() string {
	args := c.Called()
	return args.String(0)
}

func (c *ConfigMock) GetTelegramBotImageToken() string {
	args := c.Called()
	return args.String(0)
}

func (c *ConfigMock) GetGptApiKey() string {
	args := c.Called()
	return args.String(0)
}

func (c *ConfigMock) GetSendImageByUrl() bool {
	args := c.Called()
	return args.Bool(0)
}

func (c *ConfigMock) GetGptModel() string {
	args := c.Called()
	return args.String(0)
}

func (c *ConfigMock) GetTelegramWebhookToken() string {
	args := c.Called()
	return args.String(0)
}
