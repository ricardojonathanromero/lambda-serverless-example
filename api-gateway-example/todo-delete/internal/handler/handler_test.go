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
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-delete/internal/domain"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-delete/internal/handler"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
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

func (m *mockService) RemoveItem(id interface{}) error {
	args := m.Called(id)
	return args.Error(0)
}

var _ = Describe("unit tests", func() {
	Context("happy path", func() {
		var ms *mockService
		BeforeEach(func() {
			ms = new(mockService)
		})

		It("when int64 success event", func() {
			ms.On("RemoveItem", int64(1)).Return(nil)

			// constructor
			req := events.APIGatewayProxyRequest{
				Resource:       "/",
				Path:           "/v1/todos/:id",
				HTTPMethod:     http.MethodDelete,
				Headers:        map[string]string{"Content-Type": "application/json"},
				PathParameters: map[string]string{"id": "1"},
			}

			result := events.APIGatewayProxyResponse{
				StatusCode:      http.StatusAccepted,
				Headers:         map[string]string{"Content-Type": "application/json"},
				Body:            base64.StdEncoding.EncodeToString([]byte(`{"message":"removed"}`)),
				IsBase64Encoded: true,
			}

			res, err := handler.New(ms).HandleRequest(ctx, req)
			Ω(err).To(BeNil())
			Ω(res).NotTo(BeNil())
			Ω(res).To(Equal(result))
		})
		It("when primitive.ObjectID success event", func() {
			id := primitive.NewObjectID()
			ms.On("RemoveItem", id).Return(nil)

			// constructor
			req := events.APIGatewayProxyRequest{
				Resource:       "/",
				Path:           "/v1/todos/:id",
				HTTPMethod:     http.MethodDelete,
				Headers:        map[string]string{"Content-Type": "application/json"},
				PathParameters: map[string]string{"id": id.Hex()},
			}

			result := events.APIGatewayProxyResponse{
				StatusCode:      http.StatusAccepted,
				Headers:         map[string]string{"Content-Type": "application/json"},
				Body:            base64.StdEncoding.EncodeToString([]byte(`{"message":"removed"}`)),
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
				HTTPMethod:     http.MethodDelete,
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
		It("when RemoveItem generates an error", func() {
			ms.On("RemoveItem", int64(15)).Return(errors.New("no removed"))

			// constructor
			req := events.APIGatewayProxyRequest{
				Resource:       "/",
				Path:           "/v1/todos/:id",
				HTTPMethod:     http.MethodDelete,
				Headers:        map[string]string{"Content-Type": "application/json"},
				PathParameters: map[string]string{"id": "15"},
			}

			str, _ := json.Marshal(domain.NewErr("not_found", "no removed"))
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
