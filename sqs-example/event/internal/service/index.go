package service

import "github.com/ricardojonathanromero/lambda-serverless-example/sqs-example/event/internal/port"

type service struct {
}

var _ port.IService = (*service)(nil)

func New() port.IService {
	return &service{}
}
