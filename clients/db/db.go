package db

import (
	"fmt"
	"log"
	"os"

	"github.com/marlosl/gpt-telegram-bot/consts"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type CacheRepositoryInterface interface {
	SaveItem(item *string) error
	ItemExists(item string) bool
	GetItem(item string) (*string, error)
	DeleteItem(item string) error
}

type DBClient struct {
	TableName *string
	Session   *session.Session
}

type CacheRepository struct {
	DBClient
}

type CacheItem struct {
	PK string `json:"pk" dynamodbav:"PK"`
	SK string `json:"sk" dynamodbav:"SK"`
}

var MESSAGE = "MESSAGE"

func NewDBClient(tableName string, config *aws.Config) (*DBClient, error) {
	sess, err := session.NewSession(config)
	if err != nil {
		fmt.Println("Got an error creating a new session:")
		fmt.Println(err)
		return nil, err
	}

	return &DBClient{
		TableName: &tableName,
		Session:   sess,
	}, nil
}

func NewCacheRepository() (*CacheRepository, error) {
	tableName := os.Getenv(consts.CacheTable)
	dbClient, err := NewDBClient(tableName, nil)
	if err != nil {
		return nil, err
	}

	return &CacheRepository{
		*dbClient,
	}, nil
}

func (db *CacheRepository) SaveItem(item *string) error {
	svc := dynamodb.New(db.Session)

	dbItem := &CacheItem{
		PK: MESSAGE,
		SK: *item,
	}

	av, err := dynamodbattribute.MarshalMap(dbItem)
	if err != nil {
		log.Fatalf("Got error marshalling map: %s", err)
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: db.TableName,
	}

	_, err = svc.PutItem(input)
	if err != nil {
		log.Fatalf("Got error calling PutItem: %s", err)
	}

	fmt.Println("Successfully added '" + *item + "' to table " + *db.TableName)
	return nil
}

func (db *CacheRepository) ItemExists(item string) bool {
	dbItem, err := db.GetItem(item)
	if err != nil {
		fmt.Printf("Got error calling GetItem: %s\n", err)
		return false
	}

	return dbItem != nil && *dbItem != ""
}

func (db *CacheRepository) GetItem(item string) (*string, error) {
	svc := dynamodb.New(db.Session)

	sk := item
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"PK": {
				S: &MESSAGE,
			},
			"SK": {
				S: &sk,
			},
		},
		TableName: db.TableName,
	}

	result, err := svc.GetItem(input)
	if err != nil {
		fmt.Printf("Got error calling GetItem: %s\n", err)
	}

	dbItem := &CacheItem{}
	err = dynamodbattribute.UnmarshalMap(result.Item, dbItem)
	if err != nil {
		fmt.Printf("Got error unmarshalling: %s\n", err)
	}

	return &dbItem.SK, nil
}

func (db *CacheRepository) DeleteItem(item string) error {
	svc := dynamodb.New(db.Session)

	sk := item
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"PK": {
				S: &MESSAGE,
			},
			"SK": {
				S: &sk,
			},
		},
		TableName: db.TableName,
	}

	_, err := svc.DeleteItem(input)
	if err != nil {
		fmt.Printf("Got error calling DeleteItem: %v\n", err)
	}

	fmt.Println("Deleted '" + item + "' from table " + *db.TableName)
	return nil
}
