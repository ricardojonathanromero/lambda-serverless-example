package port

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-list/internal/domain"
)

type IHandler interface {
	HandleRequest(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
}

type IService interface {
	GetAll(limit, offset int64) (domain.Result, error)
}

type IRepository interface {
	CountDocuments() (int64, error)
	FindAll(limit int64, offset int64) ([]*domain.TModel, error)
}
