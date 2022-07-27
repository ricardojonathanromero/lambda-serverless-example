package port

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
)

type IHandler interface {
	// HandleRequest business logic/*
	HandleRequest(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
}

type IRepository interface {
	// Delete removes domain.TModel from db/*
	Delete(id interface{}) error
}

type IService interface {
	// RemoveItem removes item id from db/*
	RemoveItem(id interface{}) error
}
