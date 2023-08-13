package awsdeploy

import (
	"os"
	"path/filepath"

	"github.com/marlosl/gpt-telegram-bot/consts"

	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/cloudwatch"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/dynamodb"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/lambda"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/sqs"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

var (
	IamRoleLambdaExecution           pulumi.StringOutput
	IamPolicyLambdaExecution         *iam.RolePolicy
	ChatGPTTalkHandlerLambdaFunction *lambda.Function
	SendImageHandlerLambdaFunction   *lambda.Function
	ChatGPTHandlerLogGroup           *cloudwatch.LogGroup
	SendImageHandlerLogGroup         *cloudwatch.LogGroup
	CacheDynamoDbTable               *dynamodb.Table
	SQSSendImageQueue                *sqs.Queue
)

func CreateLambdaRolePolicy(ctx *pulumi.Context) error {
	role, err := iam.NewRole(ctx, "ChatGPTIamRoleLambdaExecution", &iam.RoleArgs{
		AssumeRolePolicy: pulumi.String(`{
			"Version": "2012-10-17",
			"Statement": [{
				"Sid": "",
				"Effect": "Allow",
				"Principal": {
					"Service": "lambda.amazonaws.com"
				},
				"Action": "sts:AssumeRole"
			}]
		}`),
	})
	if err != nil {
		return err
	}
	IamRoleLambdaExecution = role.Arn

	lambdaPolicy, err := iam.NewRolePolicy(ctx, "ChatGPTIamPolicyLambdaExecution", &iam.RolePolicyArgs{
		Role: role.Name,
		Policy: pulumi.String(`{
							"Version": "2012-10-17",
							"Statement": [{
									"Effect": "Allow",
									"Action": [
											"logs:CreateLogGroup",
											"logs:CreateLogStream",
											"logs:PutLogEvents"
									],
									"Resource": "arn:aws:logs:*:*:*"
							},
							{
									"Effect": "Allow",
									"Action": [
											"ssm:GetParameter"
									],
									"Resource": "arn:aws:ssm:*:*:*"
							},
							{
								"Effect": "Allow",
								"Action": [
                  	"dynamodb:PutItem",
                  	"dynamodb:UpdateItem",
                  	"dynamodb:DeleteItem",
                  	"dynamodb:BatchWriteItem",
                  	"dynamodb:GetItem",
                  	"dynamodb:BatchGetItem",
                  	"dynamodb:Scan",
                  	"dynamodb:Query"								
								],
								"Resource": "arn:aws:dynamodb:*:*:*"
							},
							{
								"Effect": "Allow",
								"Action": [
										"sqs:*"
								],
								"Resource": "arn:aws:sqs:*:*:*"
							}]
					}`),
	})

	if err != nil {
		return err
	}
	IamPolicyLambdaExecution = lambdaPolicy
	return nil
}

func CreateLambdaFunctions(ctx *pulumi.Context) error {
	outputDir := os.Getenv(consts.ProjectOutputDir)
	chatGPTTalkFile := filepath.Join(outputDir, "chatgpttalk/chat-gpt-talk-handler.zip")
	chatGPTTalkHandlerLambdaFunction, err := lambda.NewFunction(ctx, "ChatGPTTalkHandlerLambdaFunction", &lambda.FunctionArgs{
		Handler:    pulumi.String("main"),
		Role:       IamRoleLambdaExecution,
		Runtime:    pulumi.String("go1.x"),
		Name:       pulumi.String("chat-gpt-talk-handler"),
		MemorySize: pulumi.Int(128),
		Code:       pulumi.NewFileArchive(chatGPTTalkFile),
		Timeout:    pulumi.Int(300),
		Publish:    pulumi.Bool(true),
		Environment: &lambda.FunctionEnvironmentArgs{
			Variables: pulumi.StringMap{
				"REGION":           pulumi.String(os.Getenv(consts.AwsRegion)),
				"CACHE_TABLE":      CacheDynamoDbTable.Name,
				"SEND_IMAGE_QUEUE": pulumi.String("chat-gpt-send-image.fifo"),
			},
		}},
		pulumi.DependsOn([]pulumi.Resource{IamPolicyLambdaExecution, ChatGPTHandlerLogGroup, CacheDynamoDbTable}),
	)
	if err != nil {
		return err
	}
	ChatGPTTalkHandlerLambdaFunction = chatGPTTalkHandlerLambdaFunction

	sendImageHandlerFile := filepath.Join(outputDir, "chatgptsendimage/chat-gpt-send-image-handler.zip")
	sendImageHandlerLambdaFunction, err := lambda.NewFunction(ctx, "SendImageHandlerLambdaFunction", &lambda.FunctionArgs{
		Handler:    pulumi.String("main"),
		Role:       IamRoleLambdaExecution,
		Runtime:    pulumi.String("go1.x"),
		Name:       pulumi.String("chat-gpt-send-image-handler"),
		MemorySize: pulumi.Int(128),
		Code:       pulumi.NewFileArchive(sendImageHandlerFile),
		Timeout:    pulumi.Int(300),
		Publish:    pulumi.Bool(true),
		Environment: &lambda.FunctionEnvironmentArgs{
			Variables: pulumi.StringMap{
				"SEND_IMAGE_QUEUE": pulumi.String("chat-gpt-send-image.fifo"),
			},
		}},
		pulumi.DependsOn([]pulumi.Resource{IamPolicyLambdaExecution, SendImageHandlerLogGroup}),
	)
	if err != nil {
		return err
	}
	SendImageHandlerLambdaFunction = sendImageHandlerLambdaFunction

	url, err := lambda.NewFunctionUrl(ctx, "lambdaURL", &lambda.FunctionUrlArgs{
		FunctionName:      chatGPTTalkHandlerLambdaFunction.Name,
		AuthorizationType: pulumi.String("NONE"),
		Cors: &lambda.FunctionUrlCorsArgs{
			AllowOrigins: pulumi.StringArray{
				pulumi.String("*"),
			},
			AllowMethods: pulumi.StringArray{
				pulumi.String("*"),
			},
			AllowHeaders: pulumi.StringArray{
				pulumi.String("date"),
				pulumi.String("keep-alive"),
			},
			ExposeHeaders: pulumi.StringArray{
				pulumi.String("keep-alive"),
				pulumi.String("date"),
			},
			MaxAge: pulumi.Int(86400),
		},
	})
	if err != nil {
		return err
	}

	ctx.Export("url", url.FunctionUrl)
	return nil
}

func CreateLogGroups(ctx *pulumi.Context) error {
	chatGPTHandlerLogGroup, err := cloudwatch.NewLogGroup(ctx, "ChatGPTHandlerLogGroup", &cloudwatch.LogGroupArgs{
		Name: pulumi.String("/aws/lambda/chat-gpt-talk-handler"),
	})

	if err != nil {
		return err
	}

	sendImageHandlerLogGroup, err := cloudwatch.NewLogGroup(ctx, "SendImageHandlerLogGroup", &cloudwatch.LogGroupArgs{
		Name: pulumi.String("/aws/lambda/chat-gpt-send-image-handler"),
	})

	if err != nil {
		return err
	}
	ChatGPTHandlerLogGroup = chatGPTHandlerLogGroup
	SendImageHandlerLogGroup = sendImageHandlerLogGroup
	return nil
}

func CreateDynamoDbTable(ctx *pulumi.Context) error {
	cacheTable, err := dynamodb.NewTable(ctx, "CacheDynamoDbTable", &dynamodb.TableArgs{
		Attributes: dynamodb.TableAttributeArray{
			&dynamodb.TableAttributeArgs{
				Name: pulumi.String("PK"),
				Type: pulumi.String("S"),
			},
			&dynamodb.TableAttributeArgs{
				Name: pulumi.String("SK"),
				Type: pulumi.String("S"),
			},
		},
		BillingMode: pulumi.String("PAY_PER_REQUEST"),
		HashKey:     pulumi.String("PK"),
		RangeKey:    pulumi.String("SK"),
		Name:        pulumi.String("gpt-cache"),
	})
	if err != nil {
		return err
	}
	CacheDynamoDbTable = cacheTable

	return nil
}

func CreateSendImageQueue(ctx *pulumi.Context) error {
	sqsSendImageQueue, err := sqs.NewQueue(ctx, "SQSSendImageQueue", &sqs.QueueArgs{
		Name:                     pulumi.String("chat-gpt-send-image.fifo"),
		FifoQueue:                pulumi.Bool(true),
		VisibilityTimeoutSeconds: pulumi.Int(900),
	})

	if err != nil {
		return err
	}

	SQSSendImageQueue = sqsSendImageQueue
	return nil
}

func CreateLambdaEventSourceMapping(ctx *pulumi.Context) error {
	_, err := lambda.NewEventSourceMapping(ctx, "ProcessSendImageHandlerLambdaFunctionSQSSendImageEventSourceMapping", &lambda.EventSourceMappingArgs{
		EventSourceArn: SQSSendImageQueue.Arn,
		FunctionName:   SendImageHandlerLambdaFunction.Arn,
		BatchSize:      pulumi.Int(1),
		Enabled:        pulumi.Bool(true),
	},
		pulumi.DependsOn([]pulumi.Resource{IamPolicyLambdaExecution}),
	)
	if err != nil {
		return err
	}
	return nil
}
