package port

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-update/internal/domain"
)

type IHandler interface {
	// HandleRequest handle request/*
	HandleRequest(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
}

type IRepository interface {
	// Update replace current values from domain.TModel by passed/*
	Update(todo *domain.TModel) error
}

type IService interface {
	// UpdateItem updates item into db/*
	UpdateItem(item *domain.TModel) error
}
