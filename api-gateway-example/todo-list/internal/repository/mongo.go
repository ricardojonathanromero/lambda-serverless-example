package repository

import (
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-list/internal/domain"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-list/internal/port"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoRepository struct {
	conn *mongo.Client
	db   string
	col  string
}

var _ port.IRepository = (*mongoRepository)(nil)

// CountDocuments returns total documents/*
func (repo *mongoRepository) CountDocuments() (int64, error) {
	ctx, cancel := utils.NewContextWithTimeout(timeoutIn)
	defer cancel()

	return repo.collection().CountDocuments(ctx, bson.M{})
}

// FindAll returns all documents from db/*
func (repo *mongoRepository) FindAll(limit int64, offset int64) ([]*domain.TModel, error) {
	var res []*domain.TModel
	var items []*domain.MongoModel

	log.Trace("finding in db")
	opts := options.Find().SetLimit(limit).SetSkip(offset).SetSort(bson.D{{Key: "created_at", Value: 1}})

	log.Trace("init context")
	ctx, cancel := utils.NewContextWithTimeout(timeoutIn)
	defer cancel()

	result, err := repo.collection().Find(ctx, bson.M{}, opts)
	if err != nil {
		log.Errorf("documents not found: %v", err)
		return res, DocumentsNotFound
	}

	_ = result.All(ctx, &items)

	log.Tracef("document founded: %v", len(items))
	if len(items) <= 0 {
		res = make([]*domain.TModel, 0)
		return res, nil
	}

	utils.CopyStruct(items, &res)
	return res, nil
}
