package port

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-create/internal/domain"
)

type IHandler interface {
	// HandleRequest handle request/*
	HandleRequest(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
}

type IRepository interface {
	// Insert inserts a new domain.TModel document into db/*
	Insert(todo domain.TModel) (string, error)
}

type IService interface {
	// CreateTodo inserts a new item into db/*
	CreateTodo(name, priority string) (string, error)
}
