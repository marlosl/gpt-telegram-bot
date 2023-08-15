package db

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCacheRepository(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// Create a mock for the DynamoDB API
	mockDynamoDB := NewMockDynamoDBAPI(mockCtrl)

	// Create a CacheRepository instance with the mock
	repo := &CacheRepository{
		DBClient: DBClient{
			TableName: aws.String("testTable"),
			Session:   nil, // Since we're mocking DynamoDB, session isn't needed here.
		},
	}

	// Define the expected item to be used across tests
	expectedItem := "TestItem"

	t.Run("SaveItem", func(t *testing.T) {
		expectedInput := &dynamodb.PutItemInput{
			Item: map[string]*dynamodb.AttributeValue{
				"PK": {S: aws.String(MESSAGE)},
				"SK": {S: aws.String(expectedItem)},
			},
			TableName: aws.String("testTable"),
		}
		mockDynamoDB.EXPECT().PutItem(expectedInput).Return(nil, nil)
		err := repo.SaveItem(&expectedItem)
		assert.Nil(t, err)
	})

	t.Run("GetItem", func(t *testing.T) {
		expectedOutput := &dynamodb.GetItemOutput{
			Item: map[string]*dynamodb.AttributeValue{
				"PK": {S: aws.String(MESSAGE)},
				"SK": {S: aws.String(expectedItem)},
			},
		}
		mockDynamoDB.EXPECT().GetItem(gomock.Any()).Return(expectedOutput, nil)
		item, err := repo.GetItem(expectedItem)
		assert.Nil(t, err)
		assert.Equal(t, expectedItem, *item)
	})

	t.Run("ItemExists", func(t *testing.T) {
		expectedOutput := &dynamodb.GetItemOutput{
			Item: map[string]*dynamodb.AttributeValue{
				"PK": {S: aws.String(MESSAGE)},
				"SK": {S: aws.String(expectedItem)},
			},
		}
		mockDynamoDB.EXPECT().GetItem(gomock.Any()).Return(expectedOutput, nil)
		exists := repo.ItemExists(expectedItem)
		assert.True(t, exists)
	})

	t.Run("DeleteItem", func(t *testing.T) {
		mockDynamoDB.EXPECT().DeleteItem(gomock.Any()).Return(nil, nil)
		err := repo.DeleteItem(expectedItem)
		assert.Nil(t, err)
	})
}
