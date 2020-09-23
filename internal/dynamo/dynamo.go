package dynamo

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/awserr"

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
		config *Config
		svc    *dynamodb.DynamoDB
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

	return &DynamoSession{config: cfg, svc: svc}, nil
}

// Returns true when table should be created, false otherwise
func (c *DynamoSession) shouldCreateTable() (bool, error) {
	tables, err := c.DescribeTables()
	if err != nil {
		return false, err
	}

	for _, v := range tables.TableNames {
		if *v == c.config.TableName {
			return false, nil
		}
	}

	return true, nil
}

// CreateCandlesTable used to create the DynamoDB candles table
func (c *DynamoSession) CreateCandlesTable() error {
	var shouldCreate bool
	var err error
	fmt.Println(c.config.TableName)
	if shouldCreate, err = c.shouldCreateTable(); err != nil {
		return err
	}

	if !shouldCreate {
		fmt.Printf("%v table already exists, no need to make one", c.config.TableName)
		return nil
	}

	tableName := c.config.TableName

	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("Timestamp"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("Pair"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("Timestamp"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("Pair"),
				KeyType:       aws.String("RANGE"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(5),
			WriteCapacityUnits: aws.Int64(5),
		},
		TableName: aws.String(tableName),
	}

	if _, err := c.svc.CreateTable(input); err != nil {
		fmt.Println("Got error calling CreateTable")
		fmt.Println(err.Error())
		return err
	}

	fmt.Println("Created table: ", tableName)
	return nil
}

// AddItem creates an item in the table that corresponds with tableName
func (c *DynamoSession) AddItem(in interface{}) error {
	av, err := dynamodbattribute.MarshalMap(in)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(c.config.TableName),
	}

	_, err = c.svc.PutItem(input)
	if err != nil {
		return err
	}

	fmt.Printf("Successfully added %+v", in)
	return nil
}

// DescribeTables pings dynamodb and returns a list of currently active tables
func (c *DynamoSession) DescribeTables() (*dynamodb.ListTablesOutput, error) {
	input := &dynamodb.ListTablesInput{}
	result, err := c.svc.ListTables(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeInternalServerError:
				fmt.Println(dynamodb.ErrCodeInternalServerError, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			fmt.Println(err.Error())
		}
		return nil, err
	}
	fmt.Println("Active tables: ", result)
	return result, nil

}

// HasHealthyConnection ensures the dynamo connection is healthy by describing the tables
func (c *DynamoSession) HasHealthyConnection() bool {
	if _, err := c.DescribeTables(); err != nil {
		return false
	}
	fmt.Println("Connection is healthy")
	return true
}
