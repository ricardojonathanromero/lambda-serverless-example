package service_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-update/internal/domain"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-update/internal/service"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type mockRepo struct {
	mock.Mock
}

func (m *mockRepo) Update(item *domain.TModel) error {
	args := m.Called(item)
	return args.Error(0)
}

var _ = Describe("unit tests", func() {
	Context("happy path", func() {
		var mr *mockRepo
		BeforeEach(func() {
			mr = new(mockRepo)
		})

		It("return result", func() {
			input := &domain.TModel{
				ID:        primitive.NewObjectID(),
				Name:      "one",
				Status:    "created",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			mr.On("Update", input).Return(nil)
			err := service.New(mr).UpdateItem(input)
			Ω(err).To(BeNil())
		})
		It("another result", func() {
			input := &domain.TModel{
				ID:        int64(1),
				Name:      "two",
				Status:    "created",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			mr.On("Update", input).Return(nil)
			err := service.New(mr).UpdateItem(input)
			Ω(err).To(BeNil())
		})
	})
})
