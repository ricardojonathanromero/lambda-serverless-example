package repository_test

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	at "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/utils"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"io"
	"os"
	"testing"
	"time"
)

type DynamoSuite struct {
	suite.Suite
	DatabaseUrl string
}

var (
	sui          *DynamoSuite
	mongoSui     *MongoSuite
	dynamoClient *dynamodb.Client
	mongoClient  *mongo.Client
	tableName    = "todos"
)

const (
	dynamodbLocalImage string = "amazon/dynamodb-local:latest"
	timeoutIn                 = 30 * time.Second
	dbMongo                   = "todos"
	colMongo                  = "todos"
)

func (suite *DynamoSuite) StartDynamoLocal() error {
	cli, err := client.NewClientWithOpts(client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	ctx := context.Background()
	listImages, err := cli.ImageList(ctx, types.ImageListOptions{All: false})
	if err != nil {
		return err
	}

	imageExist := false
	for _, image := range listImages {
		for _, repoTag := range image.RepoTags {
			if dynamodbLocalImage == repoTag {
				imageExist = true
				break
			}
		}
		if imageExist {
			break
		}
	}

	if !imageExist {
		out, err := cli.ImagePull(ctx, dynamodbLocalImage, types.ImagePullOptions{})

		if err != nil {
			return err
		}
		_, _ = io.Copy(os.Stdout, out)
	}

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image:        dynamodbLocalImage,
		ExposedPorts: nat.PortSet{"8000": struct{}{}},
	}, &container.HostConfig{
		PortBindings: map[nat.Port][]nat.PortBinding{"8000": {{HostIP: "0.0.0.0", HostPort: "8000"}}},
	}, nil, nil, "dynamodblocal")
	if err != nil {
		return err
	}

	if err = cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}
	r, _ := cli.ContainerInspect(ctx, resp.ID)
	_ = r.Config

	if os.Getenv("CI") == "CI" {
		suite.DatabaseUrl = "http://" + r.NetworkSettings.IPAddress + ":8000"
	} else {
		suite.DatabaseUrl = "http://" + "localhost" + ":8000"
	}
	return nil
}

func (suite *DynamoSuite) CreateTable(input *dynamodb.CreateTableInput) error {
	ctx, cancel := utils.NewContextWithTimeout(timeoutIn)
	defer cancel()

	_, err1 := suite.GetDynamoDBLocalClient(ctx).CreateTable(ctx, input)
	if err1 != nil {
		return fmt.Errorf("%v\n", err1)
	}
	return nil
}

func (suite *DynamoSuite) PutItem(tableName string, item interface{}) error {
	ctx, cancel := utils.NewContextWithTimeout(timeoutIn)
	defer cancel()

	values := utils.ToDynamoDBMap(item)
	putInput := &dynamodb.PutItemInput{
		Item:      values,
		TableName: aws.String(tableName),
	}
	_, err2 := suite.GetDynamoDBLocalClient(ctx).PutItem(ctx, putInput)
	if err2 != nil {
		return fmt.Errorf("%v\n", err2)
	}
	return nil
}

func (suite *DynamoSuite) DeleteItem(tableName string, key string, value string) error {
	ctx, cancel := utils.NewContextWithTimeout(timeoutIn)
	defer cancel()

	input := &dynamodb.DeleteItemInput{
		Key: map[string]at.AttributeValue{
			key: &at.AttributeValueMemberS{Value: value},
		},
		TableName: aws.String(tableName),
	}
	_, err := suite.GetDynamoDBLocalClient(ctx).DeleteItem(ctx, input)
	if err != nil {
		return fmt.Errorf("%v\n", err)
	}
	return nil
}

func (suite *DynamoSuite) PauseContainer() {
	cli, err := client.NewClientWithOpts(client.WithAPIVersionNegotiation())
	ctx := context.Background()
	if err != nil {
		panic(err)
	}

	err = cli.ContainerPause(ctx, "dynamodblocal")
	if err != nil {
		log.Printf("container no active: %v", err)
	}
}

func (suite *DynamoSuite) ReRunContainer() {
	cli, err := client.NewClientWithOpts(client.WithAPIVersionNegotiation())
	ctx := context.Background()
	if err != nil {
		panic(err)
	}

	err = cli.ContainerUnpause(ctx, "dynamodblocal")
	if err != nil {
		log.Printf("container no active: %v", err)
	}
}

func (suite *DynamoSuite) ShutDownDynamoLocal() {
	cli, err := client.NewClientWithOpts(client.WithAPIVersionNegotiation())
	ctx := context.Background()
	if err != nil {
		panic(err)
	}

	if err = cli.ContainerStop(ctx, "dynamodblocal", nil); err != nil {
		log.Printf("Unable to stop container %s: %s", "dynamodblocal", err)
	}

	removeOptions := types.ContainerRemoveOptions{
		RemoveVolumes: true,
		Force:         true,
	}

	if err = cli.ContainerRemove(ctx, "dynamodblocal", removeOptions); err != nil {
		fmt.Printf("Unable to remove container: %v\n", err)
	}
}

func (suite *DynamoSuite) GetDynamoDBLocalClient(ctx context.Context) *dynamodb.Client {
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("us-east-1"),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{URL: suite.DatabaseUrl}, nil
			})),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID: "dummy", SecretAccessKey: "dummy", SessionToken: "dummy",
				Source: "Hard-coded credentials; values are irrelevant for local DynamoDB",
			},
		}),
	)

	if err != nil {
		log.Fatalf("error connecting to DynamoDB: %v", err)
	}

	return dynamodb.NewFromConfig(cfg)
}

type MongoSuite struct {
	suite.Suite
	mongoURI string
}

const mongoImage string = "mongo:latest"

func (suite *MongoSuite) StartMongoDBLocal() error {

	cli, err := client.NewClientWithOpts(client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	ctx := context.Background()
	listImages, err := cli.ImageList(ctx, types.ImageListOptions{All: false})
	if err != nil {
		return err
	}

	imageExist := false
	for _, image := range listImages {
		for _, repoTag := range image.RepoTags {
			if mongoImage == repoTag {
				imageExist = true
				break
			}
		}
		if imageExist {
			break
		}
	}

	if !imageExist {
		out, err := cli.ImagePull(ctx, mongoImage, types.ImagePullOptions{})

		if err != nil {
			return err
		}
		_, _ = io.Copy(os.Stdout, out)
	}

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image:        mongoImage,
		ExposedPorts: nat.PortSet{"27017": struct{}{}},
	}, &container.HostConfig{
		PortBindings: map[nat.Port][]nat.PortBinding{"27017": {{HostIP: "0.0.0.0", HostPort: "27017"}}},
	}, nil, nil, "mongolocal")
	if err != nil {
		return err
	}

	if err = cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}
	r, _ := cli.ContainerInspect(ctx, resp.ID)
	_ = r.Config

	if os.Getenv("CI") == "CI" {
		suite.mongoURI = "mongodb://" + r.NetworkSettings.IPAddress + ":27017/"
	} else {
		suite.mongoURI = "mongodb://localhost:27017/"
	}
	return nil
}

func (suite *MongoSuite) PauseContainer() {
	cli, err := client.NewClientWithOpts(client.WithAPIVersionNegotiation())
	ctx := context.Background()
	if err != nil {
		panic(err)
	}

	err = cli.ContainerPause(ctx, "mongolocal")
	if err != nil {
		log.Printf("container no active: %v", err)
	}
}

func (suite *MongoSuite) ReRunContainer() {
	cli, err := client.NewClientWithOpts(client.WithAPIVersionNegotiation())
	ctx := context.Background()
	if err != nil {
		panic(err)
	}

	err = cli.ContainerUnpause(ctx, "mongolocal")
	if err != nil {
		log.Printf("container no active: %v", err)
	}
}

func (suite *MongoSuite) ShutDownMongoLocal() {
	cli, err := client.NewClientWithOpts(client.WithAPIVersionNegotiation())
	ctx := context.Background()
	if err != nil {
		panic(err)
	}

	if err = cli.ContainerStop(ctx, "mongolocal", nil); err != nil {
		log.Printf("Unable to stop container %s: %s", "mongolocal", err)
	}

	removeOptions := types.ContainerRemoveOptions{
		RemoveVolumes: true,
		Force:         true,
	}

	if err = cli.ContainerRemove(ctx, "mongolocal", removeOptions); err != nil {
		fmt.Printf("Unable to remove container: %v\n", err)
	}
}

func (suite *MongoSuite) GetMongoDBClient() *mongo.Client {
	// create context with timeout to validate db connection
	ctx, cancel := utils.NewContextWithTimeout(timeoutIn)
	defer cancel()

	// create client connection
	log.Infof("conencting to %s", suite.mongoURI)
	mongoClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(suite.mongoURI))
	if err != nil {
		log.Fatalf("error connecting to db. reason\n%v", err)
	}

	// confirm connection making ping to db server
	log.Info("doing ping to db")
	if err = mongoClient.Ping(ctx, readpref.Primary()); err != nil {
		_ = mongoClient.Disconnect(context.Background())
		log.Fatalf("error ping to db. reason \n%v", err)
	}

	return mongoClient
}

func getTable() *dynamodb.CreateTableInput {
	return &dynamodb.CreateTableInput{
		AttributeDefinitions: []at.AttributeDefinition{
			{
				AttributeName: aws.String("Id"),
				AttributeType: at.ScalarAttributeTypeN,
			},
		},
		KeySchema: []at.KeySchemaElement{
			{
				AttributeName: aws.String("Id"),
				KeyType:       at.KeyTypeHash,
			},
		},
		GlobalSecondaryIndexes: []at.GlobalSecondaryIndex{
			{
				IndexName: aws.String("Id-Index"),
				KeySchema: []at.KeySchemaElement{
					{AttributeName: aws.String("Id"), KeyType: at.KeyTypeHash},
				},
				Projection:            &at.Projection{ProjectionType: at.ProjectionTypeAll},
				ProvisionedThroughput: &at.ProvisionedThroughput{ReadCapacityUnits: aws.Int64(5), WriteCapacityUnits: aws.Int64(5)},
			},
		},
		ProvisionedThroughput: &at.ProvisionedThroughput{ReadCapacityUnits: aws.Int64(5), WriteCapacityUnits: aws.Int64(5)},
		TableName:             aws.String(tableName),
	}
}

var _ = BeforeSuite(func() {
	ctx, cancel := utils.NewContextWithTimeout(timeoutIn)
	defer cancel()

	sui = new(DynamoSuite)
	err := sui.StartDynamoLocal()
	if err != nil {
		log.Fatalf("error starting dynamodb local: %v", err)
	}
	err = sui.CreateTable(getTable())
	if err != nil {
		sui.ShutDownDynamoLocal()
		log.Fatalf("error creating table: %v", err)
	}

	dynamoClient = sui.GetDynamoDBLocalClient(ctx)

	mongoSui = new(MongoSuite)
	err = mongoSui.StartMongoDBLocal()
	if err != nil {
		log.Fatalf("error starting mongodb local: %v", err)
	}

	mongoClient = mongoSui.GetMongoDBClient()
})

var _ = AfterSuite(func() {
	sui.ShutDownDynamoLocal()
	mongoSui.ShutDownMongoLocal()
})

func TestConn(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "repository suite test")
}
