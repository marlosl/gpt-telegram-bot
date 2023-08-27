package ssm

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

var sess *session.Session

func Get(name string) string {
	if sess == nil {
		sess = session.New()
	}
	svc := ssm.New(sess)

	output, err := svc.GetParameter(
		&ssm.GetParameterInput{
			Name:           aws.String(name),
			WithDecryption: aws.Bool(true),
		},
	)

	if err != nil {
		fmt.Printf("Error getting SSM value: %v\n", err)
		return ""
	}

	return aws.StringValue(output.Parameter.Value)
}
