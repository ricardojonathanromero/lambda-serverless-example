package port

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
)

type IHandler interface {
	// HandleEvent handle request/*
	HandleEvent(ctx context.Context, e events.SQSEvent) error
}

type IService interface {
}
