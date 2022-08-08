package handler

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/ricardojonathanromero/lambda-serverless-example/sqs-example/event/internal/port"
)

type handle struct {
	srv port.IService
}

var _ port.IHandler = (*handle)(nil)

// HandleEvent handle request/*
func (h *handle) HandleEvent(_ context.Context, _ events.SQSEvent) error {
	return nil
}

func New(srv port.IService) port.IHandler {
	return &handle{srv: srv}
}
