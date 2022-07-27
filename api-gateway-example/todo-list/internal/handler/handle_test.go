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
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-list/internal/domain"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-list/internal/handler"
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

type mockHandle struct {
	mock.Mock
}

func (m *mockHandle) GetAll(limit, offset int64) (domain.Result, error) {
	args := m.Called(limit, offset)
	return args.Get(0).(domain.Result), args.Error(1)
}

var _ = Describe("unit test", func() {
	Context("happy path", func() {
		var ms *mockHandle

		BeforeEach(func() {
			ms = new(mockHandle)
		})

		It("when handle returns happy path", func() {
			req := events.APIGatewayProxyRequest{
				Resource:              "/",
				Path:                  "/v1/todos",
				HTTPMethod:            http.MethodGet,
				Headers:               map[string]string{"Content-Type": "application/json"},
				QueryStringParameters: map[string]string{"limit": "10", "offset": "0"},
			}
			success := domain.Result{Data: getDocs(), Metadata: domain.Metadata{Limit: 10, Offset: 0, Total: 100}}
			str, _ := json.Marshal(success)

			result := events.APIGatewayProxyResponse{
				StatusCode:      http.StatusOK,
				Headers:         map[string]string{"Content-Type": "application/json"},
				Body:            base64.StdEncoding.EncodeToString(str),
				IsBase64Encoded: true,
			}

			ms.On("GetAll", int64(10), int64(0)).Return(success, nil)
			res, err := handler.New(ms).HandleRequest(ctx, req)
			Ω(err).To(BeNil())
			Ω(res).NotTo(BeNil())
			Ω(res).To(Equal(result))
		})

		It("when handle req limit is zero returns happy path", func() {
			req := events.APIGatewayProxyRequest{
				Resource:              "/",
				Path:                  "/v1/todos",
				HTTPMethod:            http.MethodGet,
				Headers:               map[string]string{"Content-Type": "application/json"},
				QueryStringParameters: map[string]string{"limit": "0", "offset": "0"},
			}
			success := domain.Result{Data: getDocs(), Metadata: domain.Metadata{Limit: 10, Offset: 0, Total: 100}}
			str, _ := json.Marshal(success)

			result := events.APIGatewayProxyResponse{
				StatusCode:      http.StatusOK,
				Headers:         map[string]string{"Content-Type": "application/json"},
				Body:            base64.StdEncoding.EncodeToString(str),
				IsBase64Encoded: true,
			}

			ms.On("GetAll", int64(10), int64(0)).Return(success, nil)
			res, err := handler.New(ms).HandleRequest(ctx, req)
			Ω(err).To(BeNil())
			Ω(res).NotTo(BeNil())
			Ω(res).To(Equal(result))
		})
	})

	Context("errors", func() {
		var ms *mockHandle

		BeforeEach(func() {
			ms = new(mockHandle)
		})

		It("when handle returns an error", func() {
			req := events.APIGatewayProxyRequest{
				Resource:              "/",
				Path:                  "/v1/todos",
				HTTPMethod:            http.MethodGet,
				Headers:               map[string]string{"Content-Type": "application/json"},
				QueryStringParameters: map[string]string{"limit": "10", "offset": "0"},
			}

			str, _ := json.Marshal(domain.NewErr("not_found", "no results"))

			result := events.APIGatewayProxyResponse{
				StatusCode:      http.StatusConflict,
				Headers:         map[string]string{"Content-Type": "application/json"},
				Body:            base64.StdEncoding.EncodeToString(str),
				IsBase64Encoded: true,
			}

			ms.On("GetAll", int64(10), int64(0)).Return(domain.Result{}, errors.New("no results"))
			res, err := handler.New(ms).HandleRequest(ctx, req)
			Ω(err).To(BeNil())
			Ω(res).NotTo(BeNil())
			Ω(res).To(Equal(result))
		})
	})
})
