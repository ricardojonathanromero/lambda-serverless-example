package domain_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-list/internal/domain"
)

var _ = Describe("unit test", func() {
	Context("create error", func() {
		It("when error is generated", func() {
			e := domain.NewErr("generic_error", "an unusual error occurred")
			Ω(e).NotTo(BeNil())
			Ω(e.Code).To(Equal("generic_error"))
			Ω(e.Message).To(Equal("an unusual error occurred"))
		})
	})
})
