//go:build wireinject
// +build wireinject
package sqs

import "github.com/google/wire"

var SQSSet = wire.NewSet(
	newSQSClient,
	wire.Bind(new(SQSClientInterface), new(*SQSClient))
)

func NewSQSClient(queue *string) *SQSClient {
  wire.Build(SQSSet)
  return SQSClient{}
}