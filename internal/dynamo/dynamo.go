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
	shouldCreateTable() (bool, error)
	CreateCandlesTable() (string, error)
	AddItem(in interface{}) error
	IsHealthy() (bool, error)
}

type (
	Config struct {
		TableName string
	}

	DynamoSession struct {
		config *Config
		ddb    *dynamodb.DynamoDB
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

	ddb := dynamodb.New(sess)

	return &DynamoSession{config: cfg, ddb: ddb}, nil
}

// Returns true when tabled does not exist and should be created, false otherwise
func (c *DynamoSession) shouldCreateTable() (bool, error) {
	input := &dynamodb.ListTablesInput{}
	output, err := c.ddb.ListTables(input)
	if err != nil {
		return false, err
	}

	for _, v := range output.TableNames {
		if *v == c.config.TableName {
			return false, nil
		}
	}

	return true, nil
}

// CreateCandlesTable used to create the DynamoDB candles table and returns a table status
func (c *DynamoSession) CreateCandlesTable() (string, error) {

	tableFailedMsg := "Create Table Failed"

	shouldCreate, err := c.shouldCreateTable()
	if err != nil {
		return tableFailedMsg, err
	}

	if shouldCreate == false {
		fmt.Printf("%+v table already exists, no need to make one\n", c.config.TableName)
		return tableFailedMsg, nil
	}

	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("Timestamp"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("PairName"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("PairName"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("Timestamp"),
				KeyType:       aws.String("RANGE"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(5),
			WriteCapacityUnits: aws.Int64(5),
		},
		TableName: aws.String(c.config.TableName),
	}

	output, err := c.ddb.CreateTable(input)
	if err != nil {
		fmt.Println("Got error calling CreateTable")
		fmt.Println(err.Error())
		return tableFailedMsg, err
	}

	fmt.Println("Created table: ", c.config.TableName)
	return *output.TableDescription.TableStatus, nil
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

	output, err := c.ddb.PutItem(input)
	if err != nil {
		return err
	}

	fmt.Println(output.GoString())
	fmt.Printf("Successfully added %+v", in)
	return nil
}

// IsHealthy returns true if the table connection is healthy, false otherwise
func (c *DynamoSession) IsHealthy() (bool, error) {
	input := &dynamodb.DescribeTableInput{TableName: &c.config.TableName}
	output, err := c.ddb.DescribeTable(input)
	if err != nil {
		return false, err
	}

	if *output.Table.TableStatus != dynamodb.TableStatusActive {
		return false, nil
	}

	fmt.Printf("Connection is healthy...\nTables Status: %+v\n", *output.Table.TableStatus)
	return true, nil
}
