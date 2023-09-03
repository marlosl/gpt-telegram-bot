//go:build wireinject
// +build wireinject
package sqs

import "github.com/google/wire"

var SQSSet = wire.NewSet(
	NewSQSClient,
	wire.Bind(new(SQSClientInterface), new(*SQSClient))
)