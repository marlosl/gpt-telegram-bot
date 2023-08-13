package awsdeploy

import (
	"context"
	"fmt"
	"os"

	"github.com/marlosl/gpt-telegram-bot/consts"

	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optdestroy"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optup"
)

var (
	AWS_DEPLOY_COMMAND  = "aws-deploy"
	AWS_DESTROY_COMMAND = "aws-destroy"

	projectName = "chat-gpt-talk"
	stackName   = "prod"
)

func ExecuteCommand(command string) {
	region := os.Getenv(consts.AwsRegion)
	ctx := context.Background()

	s, err := auto.UpsertStackInlineSource(ctx, stackName, projectName, Deploy)
	if err != nil {
		fmt.Printf("Failed initializing: %v\n", err)
		return
	}

	w := s.Workspace()

	// for inline source programs, we must manage plugins ourselves
	err = w.InstallPlugin(ctx, "aws", "v5.30.0")
	if err != nil {
		fmt.Printf("Failed to install program plugins: %v\n", err)
		os.Exit(1)
	}

	// set stack configuration specifying the AWS region to deploy
	s.SetConfig(ctx, "aws:region", auto.ConfigValue{Value: region})

	_, err = s.Refresh(ctx)
	if err != nil {
		fmt.Printf("Failed refreshing: %v\n", err)
		return
	}

	if command == AWS_DEPLOY_COMMAND {
		stdoutStreamer := optup.ProgressStreams(os.Stdout)
		_, err = s.Up(ctx, stdoutStreamer)
	}

	if command == AWS_DESTROY_COMMAND {
		stdoutStreamer := optdestroy.ProgressStreams(os.Stdout)
		_, err = s.Destroy(ctx, stdoutStreamer)
	}

	if err != nil {
		fmt.Printf("Failed updating: %v\n", err)
		return
	}
}

func GetStackOutput(key string) string {
	region := os.Getenv(consts.AwsRegion)
	ctx := context.Background()

	s, err := auto.SelectStackInlineSource(ctx, stackName, projectName, Deploy)
	if err != nil {
		fmt.Printf("Failed initializing: %v\n", err)
		return ""
	}

	w := s.Workspace()

	err = w.InstallPlugin(ctx, "aws", "v5.30.0")
	if err != nil {
		fmt.Printf("Failed to install program plugins: %v\n", err)
		return ""
	}

	s.SetConfig(ctx, "aws:region", auto.ConfigValue{Value: region})

	_, err = s.Refresh(ctx)
	if err != nil {
		fmt.Printf("Failed refreshing: %v\n", err)
		return ""
	}

	outMap, err := s.Outputs(ctx)
	if err != nil {
		fmt.Printf("Failed getting outputs: %v\n", err)
		return ""
	}
	value := outMap[key].Value
	fmt.Printf("Value: %v\n", value)

	if value != nil {
		return value.(string)
	}
	return ""
}
