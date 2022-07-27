package handler_test

import (
	"context"
	"encoding/base64"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambdacontext"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-update/internal/domain"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-update/internal/handler"
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

func (m *mockHandle) UpdateItem(item *domain.TModel) error {
	args := m.Called(item)
	return args.Error(0)
}

var _ = Describe("unit test", func() {
	Context("happy path", func() {
		var mh *mockHandle
		BeforeEach(func() {
			mh = new(mockHandle)
		})

		It("when returns success event", func() {
			now := time.Now()
			item := &domain.TModel{
				ID:        int64(1),
				Name:      "one",
				Status:    "created",
				CreatedAt: now,
				UpdatedAt: now,
			}

			req := events.APIGatewayProxyRequest{
				Resource:   "/",
				Path:       "/v1/todos",
				HTTPMethod: http.MethodPatch,
				Headers:    map[string]string{"Content-Type": "application/json"},
				Body:       utils.ToString(item),
			}

			result := events.APIGatewayProxyResponse{
				StatusCode:      http.StatusOK,
				Headers:         map[string]string{"Content-Type": "application/json"},
				Body:            base64.StdEncoding.EncodeToString([]byte(`{"message":"updated"}`)),
				IsBase64Encoded: true,
			}

			mh.On("UpdateItem", mock.Anything).Return(nil)

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

		It("when payload is null", func() {
			req := events.APIGatewayProxyRequest{
				Resource:   "/",
				Path:       "/v1/todos",
				HTTPMethod: http.MethodPatch,
				Headers:    map[string]string{"Content-Type": "application/json"},
				Body:       "",
			}

			strRes := `{"code":"invalid_request","message":"payload is not valid"}`

			result := events.APIGatewayProxyResponse{
				StatusCode:      http.StatusBadRequest,
				Headers:         map[string]string{"Content-Type": "application/json"},
				Body:            base64.StdEncoding.EncodeToString([]byte(strRes)),
				IsBase64Encoded: true,
			}

			res, err := handler.New(mh).HandleRequest(ctx, req)
			Ω(err).To(BeNil())
			Ω(res).NotTo(BeNil())
			Ω(res).To(Equal(result))
		})
		It("when request id is not valid", func() {
			now := time.Now()
			item := &domain.TModel{
				ID:        "error",
				Name:      "one",
				Status:    "created",
				CreatedAt: now,
				UpdatedAt: now,
			}

			req := events.APIGatewayProxyRequest{
				Resource:   "/",
				Path:       "/v1/todos",
				HTTPMethod: http.MethodPatch,
				Headers:    map[string]string{"Content-Type": "application/json"},
				Body:       utils.ToString(item),
			}

			strRes := `{"code":"invalid_request","message":"the id sent is not valid"}`
			result := events.APIGatewayProxyResponse{
				StatusCode:      http.StatusBadRequest,
				Headers:         map[string]string{"Content-Type": "application/json"},
				Body:            base64.StdEncoding.EncodeToString([]byte(strRes)),
				IsBase64Encoded: true,
			}

			res, err := handler.New(mh).HandleRequest(ctx, req)
			Ω(err).To(BeNil())
			Ω(res).NotTo(BeNil())
			Ω(res).To(Equal(result))
		})
		It("when update item returns an error", func() {
			now := time.Now()
			item := &domain.TModel{
				ID:        int64(1),
				Name:      "error",
				Status:    "created",
				CreatedAt: now,
				UpdatedAt: now,
			}

			req := events.APIGatewayProxyRequest{
				Resource:   "/",
				Path:       "/v1/todos",
				HTTPMethod: http.MethodPatch,
				Headers:    map[string]string{"Content-Type": "application/json"},
				Body:       utils.ToString(item),
			}

			result := events.APIGatewayProxyResponse{
				StatusCode:      http.StatusConflict,
				Headers:         map[string]string{"Content-Type": "application/json"},
				Body:            base64.StdEncoding.EncodeToString([]byte(`{"code":"not_updated","message":"no updated"}`)),
				IsBase64Encoded: true,
			}

			mh.On("UpdateItem", mock.Anything).Return(errors.New("no updated"))

			res, err := handler.New(mh).HandleRequest(ctx, req)
			Ω(err).To(BeNil())
			Ω(res).NotTo(BeNil())
			Ω(res).To(Equal(result))
		})
	})
})
