package service

import (
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-update/internal/domain"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-update/internal/port"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

type service struct {
	repo port.IRepository
}

var _ port.IService = (*service)(nil)

// UpdateItem updates item into db/*
func (s *service) UpdateItem(item *domain.TModel) error {
	log.Info("updating item...")

	return s.repo.Update(item)
}

// New constructor/*
func New(repo port.IRepository) port.IService {
	return &service{repo: repo}
}
