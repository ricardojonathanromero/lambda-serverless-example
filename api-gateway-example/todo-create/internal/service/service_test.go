package service_test

import (
	"errors"
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-create/internal/domain"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-create/internal/service"
	"github.com/stretchr/testify/mock"
	"time"
)

type mockRepo struct {
	mock.Mock
}

func (m *mockRepo) Insert(todo domain.TModel) (string, error) {
	args := m.Called(todo)
	return args.String(0), args.Error(1)
}

var _ = Describe("unit tests", func() {
	Context("happy path", func() {
		var mr *mockRepo
		BeforeEach(func() {
			mr = new(mockRepo)
		})

		It("return result", func() {
			now := time.Now()
			mr.On("Insert", mock.Anything).Return(fmt.Sprintf("%d", now.Unix()), nil)
			id, err := service.New(mr).CreateTodo("test", "URGENT")
			Ω(err).To(BeNil())
			Ω(id).NotTo(BeEmpty())
			Ω(id).To(Equal(fmt.Sprintf("%d", now.Unix())))
		})

		It("error", func() {
			mr.On("Insert", mock.Anything).Return("", errors.New("document not created"))
			id, err := service.New(mr).CreateTodo("error", "URGENT")
			Ω(id).To(BeEmpty())
			Ω(err).To(Equal(errors.New("document not created")))
		})
	})
})
