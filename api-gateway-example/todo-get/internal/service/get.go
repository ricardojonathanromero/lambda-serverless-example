package service

import (
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-get/internal/domain"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-get/internal/port"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

type service struct {
	repo port.IRepository
}

var _ port.IService = (*service)(nil)

// GetTodo looking for *domain.TModel using primitive.ObjectID/*
func (s *service) GetTodo(id interface{}) (*domain.TModel, error) {
	log.Infof("looking for %v", id)

	item, err := s.repo.FindById(id)
	if err != nil {
		log.Errorf("error retrieving for: %v", err)
		return nil, err
	}

	return item, nil
}

// New constructor/*
func New(repo port.IRepository) port.IService {
	return &service{repo: repo}
}
