package sqs

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type SQSClientInterface interface {
	SendMsg(message interface{}) error
}

type SQSClient struct {
	QueueURL *string
	Session  *session.Session
}

func NewSQSClient(queue *string) (*SQSClient, error) {
	sess, err := session.NewSession()
	if err != nil {
		fmt.Println("Got an error creating a new session:")
		fmt.Println(err)
		return nil, err
	}

	result, err := GetQueueURL(sess, queue)
	if err != nil {
		fmt.Println("Got an error getting the queue URL:")
		fmt.Println(err)
		return nil, err
	}

	return &SQSClient{
		QueueURL: result.QueueUrl,
		Session:  sess,
	}, nil
}

func (s *SQSClient) SendMsg(message interface{}) error {
	svc := sqs.New(s.Session)

	body, err := json.Marshal(message)
	if err != nil {
		fmt.Println("Got an error marshalling the message:")
		fmt.Println(err)
		return err
	}

	hashBytes := sha1.Sum(body)
	hash := fmt.Sprintf("%x", hashBytes)

	fmt.Println("HASH: ", hash)
	messageSize := len(body)

	_, err = svc.SendMessage(&sqs.SendMessageInput{
		MessageAttributes: map[string]*sqs.MessageAttributeValue{
			"MessageSize": {
				DataType:    aws.String("Number"),
				StringValue: aws.String(strconv.Itoa(messageSize)),
			},
		},
		MessageBody:            aws.String(string(body)),
		QueueUrl:               s.QueueURL,
		MessageGroupId:         aws.String("messages"),
		MessageDeduplicationId: aws.String(hash),
	})

	if err != nil {
		fmt.Println("Got an error sending the message:")
		fmt.Println(err)
		return err
	}

	fmt.Println("Sent message to queue ")
	return nil
}

func GetQueueURL(sess *session.Session, queue *string) (*sqs.GetQueueUrlOutput, error) {
	svc := sqs.New(sess)

	result, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: queue,
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}
