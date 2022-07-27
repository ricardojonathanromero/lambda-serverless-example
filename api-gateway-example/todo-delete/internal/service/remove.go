package service

import (
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-delete/internal/port"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

type service struct {
	repo port.IRepository
}

var _ port.IService = (*service)(nil)

// RemoveItem removes item id from db/*
func (s *service) RemoveItem(id interface{}) error {
	log.Infof("removing id: %v", id)

	return s.repo.Delete(id)
}

// New constructor/*
func New(repo port.IRepository) port.IService {
	return &service{repo: repo}
}
