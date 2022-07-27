package service

import (
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-list/internal/domain"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-list/internal/port"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/utils"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

type service struct {
	repo port.IRepository
}

var _ port.IService = (*service)(nil)

func (srv *service) GetAll(limit, offset int64) (domain.Result, error) {
	var resp domain.Result

	log.Info("count documents")
	total, err := srv.repo.CountDocuments()
	if err != nil {
		log.Errorf("error counting documents: %v", err)
		return resp, err
	}

	if total <= 0 {
		log.Infof("no documents founded: %v", total)
		resp = domain.Result{
			Data:     make([]*domain.TModel, 0),
			Metadata: domain.Metadata{Total: total, Limit: limit, Offset: offset},
		}
		return resp, nil
	}

	log.Info("retrieving all todos")
	result, err := srv.repo.FindAll(limit, offset)
	if err != nil {
		log.Errorf("error retrieving todos: %v", err)
		return resp, err
	}

	resp = domain.Result{
		Data:     result,
		Metadata: domain.Metadata{Total: total, Limit: limit, Offset: offset},
	}
	log.Infof("return response: %v", utils.ToString(resp))
	return resp, nil
}

func New(repo port.IRepository) port.IService {
	return &service{repo: repo}
}
