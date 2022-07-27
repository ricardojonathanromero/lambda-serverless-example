package repository_test

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-update/internal/domain"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-update/internal/port"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-update/internal/repository"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/utils"
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
			var target *domain.DynamoModel
			now := time.Now()
			item := &domain.TModel{
				ID:        int64(1),
				Name:      "one",
				Status:    "created",
				CreatedAt: now,
				UpdatedAt: now,
			}
			err := repo.Update(item)
			Ω(err).To(BeNil())

			res, err := dynamoClient.GetItem(context.Background(), &dynamodb.GetItemInput{
				Key: map[string]types.AttributeValue{
					"Id": &types.AttributeValueMemberN{Value: fmt.Sprintf("%v", int64(1))},
				},
				TableName: aws.String(tableName),
			})

			Ω(err).To(BeNil())
			Ω(res).NotTo(BeNil())
			utils.DynamoMapToInterface(res.Item, &target)
			Ω(target).NotTo(BeNil())
			Ω(target.ID).To(Equal(1))
			Ω(target.Name).To(Equal("one"))
			Ω(target.Status).To(Equal("created"))
			Ω(target.CreatedAt.Format(time.RFC3339)).To(Equal(now.Format(time.RFC3339)))
			Ω(target.UpdatedAt.Format(time.RFC3339)).To(Equal(now.Format(time.RFC3339)))
		})
	})

	Context("errors", func() {
		var repo port.IRepository
		BeforeEach(func() {
			repo = repository.NewDynamo(dynamoClient, tableName)
		})

		It("when update generates an error", func() {
			item := &domain.TModel{}
			err := repo.Update(item)
			Ω(err).NotTo(BeNil())
			Ω(err).To(Equal(repository.DocumentNotUpdated))
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
			var target *domain.MongoModel
			id := primitive.NewObjectID()
			now := time.Now()
			item := &domain.MongoModel{ID: id, Name: "one", Status: "created", CreatedAt: &now, UpdatedAt: &now}
			// insert item
			res, err := mongoClient.Database(dbMongo).Collection(colMongo).InsertOne(context.TODO(), item)
			Ω(err).To(BeNil())
			Ω(res).NotTo(BeNil())

			err = repo.Update(&domain.TModel{
				ID:        id,
				Name:      "two",
				Status:    "created",
				CreatedAt: now,
				UpdatedAt: now.Add(5 * time.Second),
			})
			Ω(err).To(BeNil())

			err = mongoClient.
				Database(dbMongo).
				Collection(colMongo).
				FindOne(context.TODO(), bson.M{"_id": id}).
				Decode(&target)

			Ω(err).To(BeNil())
			Ω(target).NotTo(BeNil())
			Ω(target.ID).To(Equal(id))
			Ω(target.Name).To(Equal("two"))
			Ω(target.Status).To(Equal("created"))
		})
	})

	Context("errors", func() {
		var repo port.IRepository
		BeforeEach(func() {
			repo = repository.NewMongo(mongoClient, dbMongo, colMongo)
		})

		It("when update returns an error", func() {
			err := repo.Update(&domain.TModel{
				ID:        primitive.NewObjectID(),
				Name:      "last",
				Status:    "created",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			})
			Ω(err).NotTo(BeNil())
			Ω(err).To(Equal(repository.DocumentNotUpdated))
		})
	})
})
