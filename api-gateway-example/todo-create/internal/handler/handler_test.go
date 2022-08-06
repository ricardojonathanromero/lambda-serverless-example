package handler_test

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambdacontext"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-create/internal/domain"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-create/internal/handler"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/utils"
	"github.com/stretchr/testify/mock"
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

type mockHandle struct {
	mock.Mock
}

func (m *mockHandle) CreateTodo(name, priority string) (string, error) {
	args := m.Called(name, priority)
	return args.String(0), args.Error(1)
}

var _ = Describe("unit test", func() {
	Context("happy path", func() {
		var mh *mockHandle
		BeforeEach(func() {
			mh = new(mockHandle)
		})

		It("when returns success event", func() {
			now := time.Now()
			item := &domain.Req{
				Name:     "one",
				Priority: "URGENT",
			}

			req := events.APIGatewayProxyRequest{
				Resource:   "/",
				Path:       "/v1/todos",
				HTTPMethod: http.MethodPost,
				Headers:    map[string]string{"Content-Type": "application/json"},
				Body:       utils.ToString(item),
			}

			result := events.APIGatewayProxyResponse{
				StatusCode:      http.StatusCreated,
				Headers:         map[string]string{"Content-Type": "application/json"},
				Body:            base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("{\"id\":\"%d\"}", now.Unix()))),
				IsBase64Encoded: true,
			}

			mh.On("CreateTodo", "one", "URGENT").Return(fmt.Sprintf("%d", now.Unix()), nil)

			res, err := handler.New(mh).HandleRequest(ctx, req)
			Ω(err).To(BeNil())
			Ω(res).NotTo(BeNil())
			Ω(res).To(Equal(result))
		})
	})

	Context("errors", func() {
		var mh *mockHandle
		BeforeEach(func() {
			mh = new(mockHandle)
		})

		It("when payload is empty", func() {
			req := events.APIGatewayProxyRequest{
				Resource:   "/",
				Path:       "/v1/todos",
				HTTPMethod: http.MethodPost,
				Headers:    map[string]string{"Content-Type": "application/json"},
				Body:       utils.ToString(`{}`),
			}

			result := events.APIGatewayProxyResponse{
				StatusCode:      http.StatusBadRequest,
				Headers:         map[string]string{"Content-Type": "application/json"},
				Body:            base64.StdEncoding.EncodeToString([]byte(`{"code":"invalid_request","message":"payload is not valid"}`)),
				IsBase64Encoded: true,
			}

			res, err := handler.New(mh).HandleRequest(ctx, req)
			Ω(err).To(BeNil())
			Ω(res).NotTo(BeNil())
			Ω(res).To(Equal(result))
		})

		It("when creation returns an error", func() {
			item := &domain.Req{
				Name:     "error_test",
				Priority: "URGENT",
			}

			req := events.APIGatewayProxyRequest{
				Resource:   "/",
				Path:       "/v1/todos",
				HTTPMethod: http.MethodPost,
				Headers:    map[string]string{"Content-Type": "application/json"},
				Body:       utils.ToString(item),
			}

			result := events.APIGatewayProxyResponse{
				StatusCode:      http.StatusNotAcceptable,
				Headers:         map[string]string{"Content-Type": "application/json"},
				Body:            base64.StdEncoding.EncodeToString([]byte(`{"code":"task_not_created","message":"error inserting item"}`)),
				IsBase64Encoded: true,
			}

			mh.On("CreateTodo", "error_test", "URGENT").Return("", errors.New("error inserting item"))

			res, err := handler.New(mh).HandleRequest(ctx, req)
			Ω(err).To(BeNil())
			Ω(res).NotTo(BeNil())
			Ω(res).To(Equal(result))
		})
	})
})
