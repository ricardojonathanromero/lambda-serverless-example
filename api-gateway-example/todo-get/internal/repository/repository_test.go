package repository_test

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-get/internal/domain"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-get/internal/port"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-get/internal/repository"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/utils"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type doc struct {
	Id        int    `json:"id" dynamodbav:"Id,omitempty"`
	Name      string `json:"name" dynamodbav:"Name,omitempty"`
	Status    string `json:"status" dynamodbav:"Status,omitempty"`
	CreatedAt string `json:"created_at,omitempty" dynamodbav:"CreatedAt,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty" dynamodbav:"UpdatedAt,omitempty"`
}

func addDynamoItems(items ...doc) []int {
	ids := make([]int, 0)

	for _, item := range items {
		err := sui.PutItem(tableName, item)
		if err != nil {
			log.Errorf("error put item: %v", err)
		} else {
			ids = append(ids, item.Id)
		}
	}

	return ids
}

func removeDynamoItems(ids ...int) {
	for _, id := range ids {
		_, err := dynamoClient.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
			Key: map[string]types.AttributeValue{
				"Id": &types.AttributeValueMemberN{Value: fmt.Sprintf("%v", id)},
			},
			TableName: aws.String(tableName),
		})
		if err != nil {
			log.Errorf("error removing item: %v", err)
		}
	}
}

var _ = Describe("dynamodb unit test", func() {
	Context("happy path", func() {
		var repo port.IRepository
		BeforeEach(func() {
			repo = repository.NewDynamo(dynamoClient, tableName)
		})

		It("when repo returns happy path", func() {
			now := time.Now().Format(time.RFC3339)
			ids := addDynamoItems(doc{Id: 1, Name: "one", Status: "created", CreatedAt: now, UpdatedAt: now})
			res, err := repo.FindById(int64(1))
			Ω(err).To(BeNil())
			Ω(res).NotTo(BeNil())
			timeNow, _ := time.Parse(time.RFC3339, now)
			Ω(res).To(Equal(&domain.TModel{ID: 1, Name: "one", Status: "created", CreatedAt: timeNow, UpdatedAt: timeNow}))
			removeDynamoItems(ids...)
		})
	})

	Context("errors", func() {
		var repo port.IRepository
		BeforeEach(func() {
			repo = repository.NewDynamo(dynamoClient, tableName)
		})

		It("when repo returns an error", func() {
			res, err := repo.FindById(int64(15))
			Ω(res).To(BeNil())
			Ω(err).NotTo(BeNil())
			Ω(err).To(Equal(repository.DocumentNotFound))
		})
	})
})

type mongoDoc struct {
	Id        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name,omitempty"`
	Status    string             `json:"status" bson:"status,omitempty"`
	CreatedAt string             `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt string             `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

func addMongoItems(items ...mongoDoc) []primitive.ObjectID {
	res := make([]primitive.ObjectID, 0)

	ctx, cancel := utils.NewContextWithTimeout(60 * time.Second)
	defer cancel()

	for _, item := range items {
		id, err := mongoClient.Database(dbMongo).Collection(colMongo).InsertOne(ctx, item)
		if err != nil {
			log.Errorf("error inserting item: %v", err)
		} else {
			res = append(res, id.InsertedID.(primitive.ObjectID))
		}
	}

	return res
}

func removeMongoItems(ids ...primitive.ObjectID) {
	ctx, cancel := utils.NewContextWithTimeout(60 * time.Second)
	defer cancel()

	for _, id := range ids {
		_, err := mongoClient.Database(dbMongo).Collection(colMongo).DeleteOne(ctx, bson.M{"_id": id})
		if err != nil {
			log.Errorf("error removing item: %v", err)
		}
	}
}

var _ = Describe("mongodb unit test", func() {
	Context("happy path", func() {
		var repo port.IRepository
		BeforeEach(func() {
			repo = repository.NewMongo(mongoClient, dbMongo, colMongo)
		})

		It("when repo returns success pathway", func() {
			now := time.Now().Format(time.RFC3339)
			ids := addMongoItems(mongoDoc{Id: primitive.NewObjectID(), Name: "one", Status: "created", CreatedAt: now, UpdatedAt: now})
			res, err := repo.FindById(ids[0])
			Ω(err).To(BeNil())
			Ω(res).NotTo(BeNil())
			Ω(res.ID).To(Equal(ids[0]))
			Ω(res.Name).To(Equal("one"))
			Ω(res.Status).To(Equal("created"))
			removeMongoItems(ids...)
		})
	})

	Context("errors", func() {
		var repo port.IRepository
		BeforeEach(func() {
			repo = repository.NewMongo(mongoClient, dbMongo, colMongo)
		})

		It("when repo returns an error", func() {
			res, err := repo.FindById(primitive.NilObjectID)
			Ω(res).To(BeNil())
			Ω(err).NotTo(BeNil())
			Ω(err).To(Equal(repository.DocumentNotFound))
		})
	})
})
