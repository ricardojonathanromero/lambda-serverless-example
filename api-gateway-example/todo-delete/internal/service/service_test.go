package service_test

import (
	"errors"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-delete/internal/service"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type mockRepo struct {
	mock.Mock
}

func (m *mockRepo) Delete(id interface{}) error {
	args := m.Called(id)
	return args.Error(0)
}

var _ = Describe("unit tests", func() {
	Context("happy path", func() {
		var mr *mockRepo

		BeforeEach(func() {
			mr = new(mockRepo)
		})

		It("when pathway returns success event", func() {
			mr.On("Delete", int64(1)).Return(nil)

			err := service.New(mr).RemoveItem(int64(1))
			Ω(err).To(BeNil())
		})

		It("when pathway returns success event primitive", func() {
			id := primitive.NewObjectID()
			mr.On("Delete", id).Return(nil)

			err := service.New(mr).RemoveItem(id)
			Ω(err).To(BeNil())
		})
	})

	Context("errors", func() {
		var mr *mockRepo

		BeforeEach(func() {
			mr = new(mockRepo)
		})

		It("when returns an error", func() {
			mr.On("Delete", int64(2)).Return(errors.New("no deleted"))
			err := service.New(mr).RemoveItem(int64(2))
			Ω(err).To(Equal(errors.New("no deleted")))
		})
		It("when returns an error", func() {
			id := primitive.NewObjectID()
			mr.On("Delete", id).Return(errors.New("no deleted"))
			err := service.New(mr).RemoveItem(id)
			Ω(err).To(Equal(errors.New("no deleted")))
		})
	})
})
