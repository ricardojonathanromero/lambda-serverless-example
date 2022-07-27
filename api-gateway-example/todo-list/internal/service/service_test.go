package service_test

import (
	"errors"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-list/internal/domain"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-list/internal/service"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type mockRepo struct {
	mock.Mock
}

func (m *mockRepo) CountDocuments() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockRepo) FindAll(limit int64, offset int64) ([]*domain.TModel, error) {
	args := m.Called(limit, offset)
	return args.Get(0).([]*domain.TModel), args.Error(1)
}

func getDocs() []*domain.TModel {
	res := []*domain.TModel{
		{
			ID:        primitive.NewObjectID(),
			Name:      "one",
			Status:    "created",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        primitive.NewObjectID(),
			Name:      "two",
			Status:    "created",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        primitive.NewObjectID(),
			Name:      "three",
			Status:    "created",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        primitive.NewObjectID(),
			Name:      "four",
			Status:    "created",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        primitive.NewObjectID(),
			Name:      "five",
			Status:    "created",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        primitive.NewObjectID(),
			Name:      "six",
			Status:    "created",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        primitive.NewObjectID(),
			Name:      "seven",
			Status:    "created",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        primitive.NewObjectID(),
			Name:      "eight",
			Status:    "created",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        primitive.NewObjectID(),
			Name:      "nine",
			Status:    "created",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        primitive.NewObjectID(),
			Name:      "ten",
			Status:    "created",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	return res
}

var _ = Describe("unit test", func() {
	Context("happy path", func() {
		var mr *mockRepo

		BeforeEach(func() {
			mr = new(mockRepo)
		})

		It("when count documents returns zero", func() {
			mr.On("CountDocuments").Return(int64(0), nil)
			res, err := service.New(mr).GetAll(10, 0)
			Ω(err).To(BeNil())
			Ω(res).NotTo(BeNil())
			Ω(res).To(Equal(domain.Result{Data: make([]*domain.TModel, 0), Metadata: domain.Metadata{Limit: 10, Offset: 0, Total: 0}}))
		})

		It("when returns success flow", func() {
			// init mocks
			items := getDocs()
			mr.On("CountDocuments").Return(int64(200), nil)
			mr.On("FindAll", int64(10), int64(0)).Return(items, nil)

			res, err := service.New(mr).GetAll(10, 0)
			Ω(err).To(BeNil())
			Ω(res).NotTo(BeNil())
			Ω(res).To(Equal(domain.Result{Data: items, Metadata: domain.Metadata{Limit: 10, Offset: 0, Total: 200}}))
		})
	})

	Context("errors", func() {
		var mr *mockRepo

		BeforeEach(func() {
			mr = new(mockRepo)
		})

		It("when count documents returns an error", func() {
			mr.On("CountDocuments").Return(int64(0), errors.New("error connection to db"))
			res, err := service.New(mr).GetAll(10, 0)
			Ω(res).To(Equal(domain.Result{}))
			Ω(err).NotTo(BeNil())
			Ω(err).To(Equal(errors.New("error connection to db")))
		})

		It("when db result return an error", func() {
			// init mocks
			mr.On("CountDocuments").Return(int64(15), nil)
			mr.On("FindAll", int64(10), int64(0)).Return([]*domain.TModel{}, errors.New("no documents found"))

			res, err := service.New(mr).GetAll(10, 0)
			Ω(res).To(Equal(domain.Result{}))
			Ω(err).NotTo(BeNil())
			Ω(err).To(Equal(errors.New("no documents found")))
		})
	})
})
