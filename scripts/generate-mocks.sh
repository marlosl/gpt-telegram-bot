#!/bin/sh
export GOPATH=$HOME/go

AWS_SDK_VERSION=$(cat go.mod | grep "github.com/aws/aws-sdk-go" | head -1 | sed -n 's/.* \(v[0-9]\+\.[0-9]\+\.[0-9]\+\)/\1/p')

mockgen -source=$GOPATH/pkg/mod/github.com/aws/aws-sdk-go@$AWS_SDK_VERSION/service/dynamodb/dynamodbiface/interface.go -destination=mocks/mock_dynamodb.go

mockgen -source=$GOPATH/pkg/mod/github.com/aws/aws-sdk-go@$AWS_SDK_VERSION/service/sqs/sqsiface/interface.go -destination=mocks/mock_sqs.go

mockgen -source=$GOPATH/pkg/mod/github.com/aws/aws-sdk-go@$AWS_SDK_VERSION/service/ssm/ssmiface/interface.go -destination=mocks/mock_ssm.go
