package database

import (
	"fmt"
	"lambda-func/types"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const (
	TABLE_NAME = "userTable"
)

type UserStore interface {
	DoesUserExist(username string) (bool, error)
	InsertUser(user types.User) error
	GetUser(username string) (types.User, error)
}

type DynamoDBClient struct {
	databaseStore *dynamodb.DynamoDB
}

func NewDynamoDBClient() DynamoDBClient {
	dbSession := session.Must(session.NewSession())
	db := dynamodb.New(dbSession)
	return DynamoDBClient{
		databaseStore: db,
	}
}

func (c DynamoDBClient) DoesUserExist(username string) (bool, error) {
	result, err := c.databaseStore.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(TABLE_NAME),
		Key: map[string]*dynamodb.AttributeValue{
			"username": {
				S: aws.String(username),
			},
		},
	})

	if err != nil {
		return false, err
	}

	if result.Item == nil {
		return false, nil
	}

	return true, nil
}

func (c DynamoDBClient) InsertUser(user types.User) error {
	item := &dynamodb.PutItemInput{
		TableName: aws.String(TABLE_NAME),
		Item: map[string]*dynamodb.AttributeValue{
			"username": {
				S: aws.String(user.Username),
			},
			"password": {
				S: aws.String(user.PasswordHash),
			},
		},
	}

	_, err := c.databaseStore.PutItem(item)
	if err != nil {
		return err
	}

	return nil
}

func (c DynamoDBClient) GetUser(username string) (types.User, error) {
	user, err := c.databaseStore.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(TABLE_NAME),
		Key: map[string]*dynamodb.AttributeValue{
			"username": {
				S: aws.String(username),
			},
		},
	})

	if err != nil {
		return types.User{}, err
	}

	if user.Item == nil {
		return types.User{}, fmt.Errorf("user does not exist")
	}

	userStruct := types.User{}
	err = dynamodbattribute.UnmarshalMap(user.Item, &userStruct)
	if err != nil {
		return types.User{}, err
	}

	return userStruct, nil
}
