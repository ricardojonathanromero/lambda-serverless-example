package repository_test

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-delete/internal/port"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-delete/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

var _ = Describe("dynamo unit tests", func() {
	Context("happy path", func() {
		var repo port.IRepository
		BeforeEach(func() {
			repo = repository.NewDynamo(dynamoClient, tableName)
		})

		It("when Update returns success event", func() {
			// insert item
			_, err := dynamoClient.PutItem(context.Background(), &dynamodb.PutItemInput{
				TableName: aws.String(tableName),
				Item: map[string]types.AttributeValue{
					"Id":        &types.AttributeValueMemberN{Value: fmt.Sprintf("%v", int64(1))},
					"Name":      &types.AttributeValueMemberS{Value: "one"},
					"Status":    &types.AttributeValueMemberS{Value: "created"},
					"CreatedAt": &types.AttributeValueMemberS{Value: time.Now().Format(time.RFC3339)},
					"UpdatedAt": &types.AttributeValueMemberS{Value: time.Now().Format(time.RFC3339)},
				},
			})
			Ω(err).To(BeNil())

			err = repo.Delete(int64(1))
			Ω(err).To(BeNil())
		})
	})

	Context("errors", func() {
		var repo port.IRepository
		BeforeEach(func() {
			repo = repository.NewDynamo(dynamoClient, tableName)
		})

		It("when update generates an error", func() {
			sui.PauseContainer()
			err := repo.Delete(int64(5))
			Ω(err).To(Equal(repository.DocumentNotDeleted))
			sui.ReRunContainer()
		})
	})
})

var _ = Describe("mongo unit tests", func() {
	Context("happy path", func() {
		var repo port.IRepository
		BeforeEach(func() {
			repo = repository.NewMongo(mongoClient, dbMongo, colMongo)
		})

		It("when update item return success event", func() {
			id := primitive.NewObjectID()

			_, _ = mongoClient.Database(dbMongo).Collection(colMongo).InsertOne(context.Background(), bson.M{
				"_id":        id,
				"name":       "one",
				"status":     "created",
				"created_at": time.Now(),
				"updated_at": time.Now(),
			})

			err := repo.Delete(id)
			Ω(err).To(BeNil())
		})
	})

	Context("errors", func() {
		var repo port.IRepository
		BeforeEach(func() {
			repo = repository.NewMongo(mongoClient, dbMongo, colMongo)
		})

		It("when update returns an error", func() {
			err := repo.Delete(primitive.NewObjectID())
			Ω(err).To(Equal(repository.DocumentNotDeleted))
		})
	})
})
