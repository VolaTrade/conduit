package dynamo

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/wire"
)

var Module = wire.NewSet(
	New,
)

type Dynamo interface {
}

type (
	Config struct {
		TableName string
	}

	DynamoSession struct {
		config       *Config
		dynamoClient *dynamodb.DynamoDB
	}
)

// New creates a new DynamoSession Struct and initiates an aws session with credentials
func New(cfg *Config) (*DynamoSession, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-1"),
	})

	if err != nil {
		return nil, err
	}

	svc := dynamodb.New(sess)

	return &DynamoSession{config: cfg, dynamoClient: svc}, nil
}

// AddItem creates an item in the table that corresponds with tableName
func (c *DynamoSession) AddItem(in interface{}, tableName string) error {
	av, err := dynamodbattribute.MarshalMap(in)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = c.dynamoClient.PutItem(input)
	if err != nil {
		return err
	}

	fmt.Printf("Successfully added %+v", in)
	return nil
}
