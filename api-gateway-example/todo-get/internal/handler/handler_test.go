package handler_test

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambdacontext"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-get/internal/domain"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-get/internal/handler"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

var ct = &lambdacontext.LambdaContext{
	AwsRequestID:       "awsRequestId1234",
	InvokedFunctionArn: "arn:aws:lambda:xxx",
	Identity:           lambdacontext.CognitoIdentity{},
	ClientContext:      lambdacontext.ClientContext{},
}

var ctx = lambdacontext.NewContext(context.TODO(), ct)

type mockService struct {
	mock.Mock
}

func (m *mockService) GetTodo(id interface{}) (*domain.TModel, error) {
	args := m.Called(id)
	return args.Get(0).(*domain.TModel), args.Error(1)
}

var _ = Describe("units test", func() {
	Context("happy path", func() {
		var ms *mockService
		BeforeEach(func() {
			ms = new(mockService)
		})

		It("when GetTodo with int64 returns success event", func() {
			now := time.Now()
			input := &domain.TModel{ID: int64(1), Name: "one", Status: "created", CreatedAt: now, UpdatedAt: now}
			ms.On("GetTodo", int64(1)).Return(input, nil)

			// constructor
			req := events.APIGatewayProxyRequest{
				Resource:       "/",
				Path:           "/v1/todos/:id",
				HTTPMethod:     http.MethodGet,
				Headers:        map[string]string{"Content-Type": "application/json"},
				PathParameters: map[string]string{"id": "1"},
			}
			success := &domain.TModel{
				ID:        int64(1),
				Name:      "one",
				Status:    "created",
				CreatedAt: now,
				UpdatedAt: now,
			}
			str, _ := json.Marshal(success)

			result := events.APIGatewayProxyResponse{
				StatusCode:      http.StatusOK,
				Headers:         map[string]string{"Content-Type": "application/json"},
				Body:            base64.StdEncoding.EncodeToString(str),
				IsBase64Encoded: true,
			}

			res, err := handler.New(ms).HandleRequest(ctx, req)
			Ω(err).To(BeNil())
			Ω(res).NotTo(BeNil())
			Ω(res).To(Equal(result))
		})

		It("when GetTodo with primitive.ObjectID returns success event", func() {
			now := time.Now()
			id := primitive.NewObjectID()
			input := &domain.TModel{ID: id, Name: "one", Status: "created", CreatedAt: now, UpdatedAt: now}
			ms.On("GetTodo", id).Return(input, nil)

			// constructor
			req := events.APIGatewayProxyRequest{
				Resource:       "/",
				Path:           "/v1/todos/:id",
				HTTPMethod:     http.MethodGet,
				Headers:        map[string]string{"Content-Type": "application/json"},
				PathParameters: map[string]string{"id": id.Hex()},
			}
			success := &domain.TModel{
				ID:        id,
				Name:      "one",
				Status:    "created",
				CreatedAt: now,
				UpdatedAt: now,
			}
			str, _ := json.Marshal(success)

			result := events.APIGatewayProxyResponse{
				StatusCode:      http.StatusOK,
				Headers:         map[string]string{"Content-Type": "application/json"},
				Body:            base64.StdEncoding.EncodeToString(str),
				IsBase64Encoded: true,
			}

			res, err := handler.New(ms).HandleRequest(ctx, req)
			Ω(err).To(BeNil())
			Ω(res).NotTo(BeNil())
			Ω(res).To(Equal(result))
		})
	})

	Context("errors", func() {
		var ms *mockService
		BeforeEach(func() {
			ms = new(mockService)
		})

		It("when id is zero", func() {
			// constructor
			req := events.APIGatewayProxyRequest{
				Resource:       "/",
				Path:           "/v1/todos/:id",
				HTTPMethod:     http.MethodGet,
				Headers:        map[string]string{"Content-Type": "application/json"},
				PathParameters: map[string]string{"id": "0"},
			}

			str, _ := json.Marshal(domain.NewErr("invalid_request", "the id sent is not valid"))
			result := events.APIGatewayProxyResponse{
				StatusCode:      http.StatusBadRequest,
				Headers:         map[string]string{"Content-Type": "application/json"},
				Body:            base64.StdEncoding.EncodeToString(str),
				IsBase64Encoded: true,
			}

			res, err := handler.New(ms).HandleRequest(ctx, req)
			Ω(err).To(BeNil())
			Ω(res).NotTo(BeNil())
			Ω(res).To(Equal(result))
		})

		It("when GetTodo generates an error", func() {
			ms.On("GetTodo", int64(2)).Return(&domain.TModel{}, errors.New("no document"))

			// constructor
			req := events.APIGatewayProxyRequest{
				Resource:       "/",
				Path:           "/v1/todos/:id",
				HTTPMethod:     http.MethodGet,
				Headers:        map[string]string{"Content-Type": "application/json"},
				PathParameters: map[string]string{"id": "2"},
			}

			str, _ := json.Marshal(domain.NewErr("not_found", "no document"))
			result := events.APIGatewayProxyResponse{
				StatusCode:      http.StatusNotFound,
				Headers:         map[string]string{"Content-Type": "application/json"},
				Body:            base64.StdEncoding.EncodeToString(str),
				IsBase64Encoded: true,
			}

			res, err := handler.New(ms).HandleRequest(ctx, req)
			Ω(err).To(BeNil())
			Ω(res).NotTo(BeNil())
			Ω(res).To(Equal(result))
		})
	})
})
