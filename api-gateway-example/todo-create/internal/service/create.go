package service

import (
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-create/internal/domain"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-create/internal/port"
	"github.com/sirupsen/logrus"
	"time"
)

type service struct {
	repo port.IRepository
}

var (
	log               = logrus.New()
	_   port.IService = (*service)(nil)
)

// CreateTodo inserts a new item into db/*
func (srv *service) CreateTodo(name, priority string) (string, error) {
	var result string
	log.Info("creating todo task")

	todo := domain.TModel{
		Name:      name,
		Priority:  priority,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	id, err := srv.repo.Insert(todo)
	if err != nil {
		log.Errorf("error creating todo task: %v", err)
		return result, err
	}

	log.Info("todo created")
	return id, nil
}

func New(repo port.IRepository) port.IService {
	return &service{repo: repo}
}
