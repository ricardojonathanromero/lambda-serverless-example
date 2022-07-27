package service_test

import (
	"errors"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-get/internal/domain"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-get/internal/service"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type mockService struct {
	mock.Mock
}

func (m *mockService) FindById(id interface{}) (*domain.TModel, error) {
	args := m.Called(id)
	return args.Get(0).(*domain.TModel), args.Error(1)
}

var _ = Describe("unit tests", func() {
	Context("happy path", func() {
		var ms *mockService

		BeforeEach(func() {
			ms = new(mockService)
		})

		It("when FindById returns a success event", func() {
			now := time.Now()
			item := &domain.TModel{ID: int64(1), Name: "one", Status: "created", CreatedAt: now, UpdatedAt: now}
			ms.On("FindById", int64(1)).Return(item, nil)
			item, err := service.New(ms).GetTodo(int64(1))
			Ω(err).To(BeNil())
			Ω(item).NotTo(BeNil())
			Ω(item).To(Equal(&domain.TModel{ID: int64(1), Name: "one", Status: "created", CreatedAt: now, UpdatedAt: now}))
		})
	})

	Context("errors", func() {
		var ms *mockService

		BeforeEach(func() {
			ms = new(mockService)
		})

		It("when FindById returns an error", func() {
			id := primitive.NewObjectID()
			ms.On("FindById", id).Return(&domain.TModel{}, errors.New("no document"))
			item, err := service.New(ms).GetTodo(id)
			Ω(item).To(BeNil())
			Ω(err).NotTo(BeNil())
			Ω(err).To(Equal(errors.New("no document")))
		})
	})
})
