package port

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-get/internal/domain"
)

type IHandle interface {
	// HandleRequest business logic/*
	HandleRequest(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
}

type IRepository interface {
	// FindById returns document type domain.TModel from db/*
	FindById(id interface{}) (*domain.TModel, error)
}

type IService interface {
	// GetTodo looking for *domain.TModel using id as interface{}/*
	GetTodo(id interface{}) (*domain.TModel, error)
}
