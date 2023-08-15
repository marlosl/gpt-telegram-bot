package ssm

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// Create a mock for the SSM API
	mockSSM := NewMockSSMAPI(mockCtrl)

	// Define the expected parameter name and value
	paramName := "testParam"
	expectedValue := "testValue"

	t.Run("Get_Success", func(t *testing.T) {
		mockSSM.EXPECT().GetParameter(&ssm.GetParameterInput{
			Name:           aws.String(paramName),
			WithDecryption: aws.Bool(true),
		}).Return(&ssm.GetParameterOutput{
			Parameter: &ssm.Parameter{
				Value: aws.String(expectedValue),
			},
		}, nil)

		// Override global sess variable with mocked session
		sess = session.Must(session.NewSession()) // Mocked session

		value := Get(paramName)
		assert.Equal(t, expectedValue, value)
	})

	t.Run("Get_Error", func(t *testing.T) {
		mockSSM.EXPECT().GetParameter(&ssm.GetParameterInput{
			Name:           aws.String(paramName),
			WithDecryption: aws.Bool(true),
		}).Return(nil, errors.New("an error occurred"))

		// Override global sess variable with mocked session
		sess = session.Must(session.NewSession()) // Mocked session

		value := Get(paramName)
		assert.Equal(t, "", value)
	})
}
