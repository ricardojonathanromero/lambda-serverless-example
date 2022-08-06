package repository_test

import (
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-create/internal/domain"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-create/internal/port"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-create/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

var _ = Describe("dynamo unit tests", func() {
	Context("happy path", func() {
		var repo port.IRepository
		BeforeEach(func() {
			repo = repository.NewDynamo(dynamoClient, tableName)
		})

		It("when insert returns success event", func() {
			now := time.Now()
			todo := domain.TModel{
				ID:        int(now.Unix()),
				Name:      "example",
				Priority:  "URGENT",
				CreatedAt: now,
				UpdatedAt: now,
			}

			id, err := repo.Insert(todo)
			Ω(err).To(BeNil())
			Ω(id).NotTo(BeEmpty())
			Ω(id).To(Equal(fmt.Sprintf("%d", int(now.Unix()))))
		})
	})

	Context("errors", func() {
		var repo port.IRepository
		BeforeEach(func() {
			repo = repository.NewDynamo(dynamoClient, tableName)
		})

		It("when insert generates an error", func() {
			item := domain.TModel{}
			_, err := repo.Insert(item)
			Ω(err).NotTo(BeNil())
			Ω(err).To(Equal(repository.DocumentNotInserted))
		})
	})
})

var _ = Describe("mongo unit tests", func() {
	Context("happy path", func() {
		var repo port.IRepository
		BeforeEach(func() {
			repo = repository.NewMongo(mongoClient, dbMongo, colMongo)
		})

		It("when insert returns success event", func() {
			now := time.Now()
			todo := domain.TModel{
				Name:      "example",
				Priority:  "URGENT",
				CreatedAt: now,
				UpdatedAt: now,
			}

			id, err := repo.Insert(todo)
			Ω(err).To(BeNil())
			Ω(id).NotTo(BeEmpty())
		})
	})

	Context("errors", func() {
		var repo port.IRepository
		BeforeEach(func() {
			repo = repository.NewMongo(mongoClient, dbMongo, colMongo)
		})

		It("when insert returns an error", func() {
			mongoSui.PauseContainer()
			id, err := repo.Insert(domain.TModel{})
			Ω(id).To(Equal(primitive.NilObjectID.Hex()))
			Ω(err).NotTo(BeNil())
			Ω(err).To(Equal(repository.DocumentNotInserted))
			mongoSui.ReRunContainer()
		})
	})
})
