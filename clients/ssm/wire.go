//go:build wireinject
// +build wireinject

package ssm

import "github.com/google/wire"

var SSMSet = wire.NewSet(
	NewSSMClient,
	wire.Bind(new(SSMClientInterface), new(*SSMClient))
)
