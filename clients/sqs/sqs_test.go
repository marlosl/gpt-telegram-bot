package sqs

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestSQSClient(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// Create a mock for the SQS API
	mockSQS := NewMockSQSAPI(mockCtrl)

	queueName := "testQueue"
	expectedQueueURL := "http://localhost:4566/000000000000/testQueue"

	t.Run("NewSQSClient", func(t *testing.T) {
		mockSQS.EXPECT().GetQueueUrl(&sqs.GetQueueUrlInput{
			QueueName: &queueName,
		}).Return(&sqs.GetQueueUrlOutput{
			QueueUrl: &expectedQueueURL,
		}, nil)

		client, err := NewSQSClient(&queueName)
		assert.Nil(t, err)
		assert.Equal(t, expectedQueueURL, *client.QueueURL)
	})

	t.Run("SendMsg", func(t *testing.T) {
		message := map[string]string{
			"key": "value",
		}
		mockSQS.EXPECT().SendMessage(gomock.Any()).Return(&sqs.SendMessageOutput{}, nil)

		client := &SQSClient{
			QueueURL: &expectedQueueURL,
			Session:  nil, // As we're mocking SQS, session isn't needed here.
		}
		err := client.SendMsg(message)
		assert.Nil(t, err)
	})

	t.Run("GetQueueURL", func(t *testing.T) {
		mockSQS.EXPECT().GetQueueUrl(&sqs.GetQueueUrlInput{
			QueueName: &queueName,
		}).Return(&sqs.GetQueueUrlOutput{
			QueueUrl: &expectedQueueURL,
		}, nil)

		sess := session.Must(session.NewSession()) // Mocked session
		queueURL, err := GetQueueURL(sess, &queueName)
		assert.Nil(t, err)
		assert.Equal(t, expectedQueueURL, *queueURL.QueueUrl)
	})
}
