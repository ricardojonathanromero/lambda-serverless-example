package repository_test

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-list/internal/port"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-list/internal/repository"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"sort"
	"time"
)

type doc struct {
	Id        int    `json:"id" dynamodbav:"Id,omitempty"`
	Name      string `json:"name" dynamodbav:"Name,omitempty"`
	Status    string `json:"status" dynamodbav:"Status,omitempty"`
	CreatedAt string `json:"created_at,omitempty" dynamodbav:"CreatedAt,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty" dynamodbav:"UpdatedAt,omitempty"`
}

func addItem() doc {
	// insert object
	now := time.Now().Format(time.RFC3339)
	d := doc{
		Id:        1,
		Name:      "First task",
		Status:    "created",
		CreatedAt: now,
		UpdatedAt: now,
	}

	_ = sui.PutItem(tableName, d)
	return d
}

func addItems() []doc {
	// insert object
	inputs := []doc{
		{
			Id:        1,
			Name:      "First",
			Status:    "created",
			CreatedAt: "2022-07-24T01:01:53-05:00",
			UpdatedAt: "2022-07-24T01:01:53-05:00",
		},
		{
			Id:        2,
			Name:      "Second",
			Status:    "created",
			CreatedAt: "2022-07-24T01:01:54-05:00",
			UpdatedAt: "2022-07-24T01:01:54-05:00",
		},
		{
			Id:        3,
			Name:      "Third",
			Status:    "created",
			CreatedAt: "2022-07-24T01:01:55-05:00",
			UpdatedAt: "2022-07-24T01:01:55-05:00",
		},
		{
			Id:        4,
			Name:      "Fourth",
			Status:    "created",
			CreatedAt: "2022-07-24T01:01:56-05:00",
			UpdatedAt: "2022-07-24T01:01:56-05:00",
		},
		{
			Id:        5,
			Name:      "Fifth",
			Status:    "created",
			CreatedAt: "2022-07-24T01:01:57-05:00",
			UpdatedAt: "2022-07-24T01:01:57-05:00",
		},
		{
			Id:        6,
			Name:      "Sixth",
			Status:    "created",
			CreatedAt: "2022-07-24T01:01:58-05:00",
			UpdatedAt: "2022-07-24T01:01:58-05:00",
		},
		{
			Id:        7,
			Name:      "Seventh",
			Status:    "created",
			CreatedAt: "2022-07-24T01:01:59-05:00",
			UpdatedAt: "2022-07-24T01:01:59-05:00",
		},
		{
			Id:        8,
			Name:      "Eighth",
			Status:    "created",
			CreatedAt: "2022-07-24T01:02:00-05:00",
			UpdatedAt: "2022-07-24T01:02:00-05:00",
		},
		{
			Id:        9,
			Name:      "Ninth",
			Status:    "created",
			CreatedAt: "2022-07-24T01:02:01-05:00",
			UpdatedAt: "2022-07-24T01:02:01-05:00",
		},
		{
			Id:        10,
			Name:      "Tenth",
			Status:    "created",
			CreatedAt: "2022-07-24T01:02:02-05:00",
			UpdatedAt: "2022-07-24T01:02:02-05:00",
		},
		{
			Id:        11,
			Name:      "Eleventh",
			Status:    "created",
			CreatedAt: "2022-07-24T01:02:03-05:00",
			UpdatedAt: "2022-07-24T01:02:03-05:00",
		},
		{
			Id:        12,
			Name:      "Twelve",
			Status:    "created",
			CreatedAt: "2022-07-24T01:02:04-05:00",
			UpdatedAt: "2022-07-24T01:02:04-05:00",
		},
	}

	sort.Slice(inputs, func(i, j int) bool {
		return inputs[i].CreatedAt < inputs[j].CreatedAt
	})

	for _, input := range inputs {
		err := sui.PutItem(tableName, input)
		if err != nil {
			log.Warn(err)
		}
	}

	return inputs
}

func removeItems(docs ...doc) {
	for _, d := range docs {
		_, err := dynamoClient.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
			Key: map[string]types.AttributeValue{
				"Id": &types.AttributeValueMemberN{Value: fmt.Sprintf("%v", d.Id)},
			},
			TableName: aws.String(tableName),
		})
		if err != nil {
			log.Errorf("error removing item: %v", err)
		}
	}
}

var _ = Describe("dynamodb repo unit test", func() {
	Context("happy path", func() {
		var repo port.IRepository
		BeforeEach(func() {
			repo = repository.NewDynamoRepo(dynamoClient, tableName)
		})
		It("when dynamodb count return zero elements", func() {
			total, err := repo.CountDocuments()
			Ω(err).To(BeNil())
			Ω(total).To(Equal(int64(0)))
		})

		It("when count return one element", func() {
			id := addItem()
			total, err := repo.CountDocuments()
			Ω(err).To(BeNil())
			Ω(total).To(Equal(int64(1)))
			removeItems(id)
		})

		It("when FindAll response successfully with zero item", func() {
			res, err := repo.FindAll(10, 0)
			Ω(err).To(BeNil())
			Ω(res).NotTo(BeNil())
			Ω(len(res)).To(Equal(0))
		})

		It("when FindAll response successfully with one item", func() {
			id := addItem()
			time.Sleep(1 * time.Second)
			res, err := repo.FindAll(10, 0)
			Ω(err).To(BeNil())
			Ω(res).NotTo(BeNil())
			Ω(len(res)).To(Equal(1))
			Ω(res[0].Name).To(Equal("First task"))
			removeItems(id)
		})

		It("when FindAll response successfully with various items", func() {
			lm := int64(10)
			docs := addItems()
			pointer := docs[lm-1].Id
			res, err := repo.FindAll(lm, int64(pointer))
			Ω(err).To(BeNil())
			Ω(res).NotTo(BeNil())
			Ω(len(res)).To(Equal(2))
			removeItems(docs...)
		})
	})

	Context("errors", func() {
		var repo port.IRepository
		BeforeEach(func() {
			repo = repository.NewDynamoRepo(dynamoClient, tableName)
		})
		It("when FindAll response an error", func() {
			lm := int64(10)
			item := addItem()
			sui.PauseContainer()
			res, err := repo.FindAll(lm, 1)
			Ω(res).To(BeNil())
			Ω(err).NotTo(BeNil())
			Ω(err).To(Equal(repository.DocumentsNotFound))
			sui.ReRunContainer()
			removeItems(item)
		})

		It("when client returns an error", func() {
			sui.PauseContainer()
			total, err := repo.CountDocuments()
			Ω(err).NotTo(BeNil())
			Ω(total).To(Equal(int64(0)))
			sui.ReRunContainer()
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

func addMongoItems(inputs ...mongoDoc) []primitive.ObjectID {
	ids := make([]primitive.ObjectID, 0)

	for _, input := range inputs {
		id := primitive.NewObjectID()
		input.Id = id

		_, err := mongoClient.Database(dbMongo).Collection(colMongo).InsertOne(context.Background(), input)
		if err != nil {
			log.Errorf("error insertign doc: %v", err)
		} else {
			ids = append(ids, id)
		}
	}

	return ids
}

func removeMongoItems(ids ...primitive.ObjectID) {
	for _, id := range ids {
		_, err := mongoClient.Database(dbMongo).Collection(colMongo).DeleteOne(context.Background(), bson.M{"_id": id})
		if err != nil {
			log.Errorf("error deleting id: %v", err)
		}
	}
}

var _ = Describe("mongodb repo unit test", func() {
	Context("happy path", func() {
		var repo port.IRepository
		BeforeEach(func() {
			repo = repository.NewMongoRepo(mongoClient, dbMongo, colMongo)
		})

		It("when dynamodb count return zero elements", func() {
			total, err := repo.CountDocuments()
			Ω(err).To(BeNil())
			Ω(total).To(Equal(int64(0)))
		})

		It("when count return one element", func() {
			now := time.Now().Format(time.RFC3339)
			id := addMongoItems(mongoDoc{Name: "one", Status: "created", CreatedAt: now, UpdatedAt: now})
			total, err := repo.CountDocuments()
			Ω(err).To(BeNil())
			Ω(total).To(Equal(int64(1)))
			removeMongoItems(id...)
		})

		It("when FindAll response successfully with zero item", func() {
			res, err := repo.FindAll(10, 0)
			Ω(err).To(BeNil())
			Ω(res).NotTo(BeNil())
			Ω(len(res)).To(Equal(0))
		})

		It("when FindAll response successfully with one item", func() {
			now := time.Now().Format(time.RFC3339)
			id := addMongoItems(mongoDoc{Name: "one", Status: "created", CreatedAt: now, UpdatedAt: now})
			time.Sleep(1 * time.Second)
			res, err := repo.FindAll(10, 0)
			Ω(err).To(BeNil())
			Ω(res).NotTo(BeNil())
			Ω(len(res)).To(Equal(1))
			Ω(res[0].Name).To(Equal("one"))
			removeMongoItems(id...)
		})
	})

	Context("errors", func() {
		var repo port.IRepository
		BeforeEach(func() {
			repo = repository.NewMongoRepo(mongoClient, dbMongo, colMongo)
		})
		It("when FindAll returns an error", func() {
			mongoSui.PauseContainer()
			res, err := repo.FindAll(10, 2)
			mongoSui.ReRunContainer()
			Ω(res).To(BeNil())
			Ω(err).NotTo(BeNil())
			Ω(err).To(Equal(repository.DocumentsNotFound))
		})
	})
})
