package repository

import (
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-get/internal/domain"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-get/internal/port"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoRepo struct {
	conn *mongo.Client
	db   string
	col  string
}

var _ port.IRepository = (*mongoRepo)(nil)

// FindById returns document type domain.TModel from db/*
func (repo *mongoRepo) FindById(id interface{}) (*domain.TModel, error) {
	var res *domain.TModel
	var item *domain.MongoModel

	ctx, cancel := utils.NewContextWithTimeout(timeoutIn)
	defer cancel()

	err := repo.collection().FindOne(ctx, bson.M{"_id": id.(primitive.ObjectID)}).Decode(&item)
	if err != nil {
		log.Errorf("document %v not found: %v", id, err)
		return res, DocumentNotFound
	}

	res = &domain.TModel{
		ID:        item.ID,
		Name:      item.Name,
		Status:    item.Status,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}
	return res, nil
}
