//go:build wireinject
// +build wireinject

package config

import (
	"github.com/marlosl/gpt-telegram-bot/clients/ssm"

	"github.com/google/wire"
)

var ConfigSet = wire.NewSet(
	ssm.SSMSet,
	NewSSMConfig,
	wire.Bind(new(ConfigInterface), new(*Config)),
)

func InitConfig() *Config {
	wire.Build(ConfigSet)
	return &Config{}
}
