package ssm

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

type SSMClientInterface interface {
  Get(name string) string
}

type SSMClient struct {
  sess *session.Session
  svc *ssm.SSM
}

func NewSSMClient() *SSMClient {
  s := new(SSMClient)
  s.init()
  return s
}

func (s *SSMClient) init() {
  t.sess = session.New()
  t.svc = ssm.New(sess)
}

func (s *SSMClient) Get(name string) string {
	output, err := s.svc.GetParameter(
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
