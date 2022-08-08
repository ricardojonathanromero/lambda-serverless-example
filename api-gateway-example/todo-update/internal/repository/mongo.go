package repository

import (
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-update/internal/domain"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-update/internal/port"
	"github.com/ricardojonathanromero/lambda-utilities/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoRepo struct {
	conn *mongo.Client
	db   string
	col  string
}

var _ port.IRepository = (*mongoRepo)(nil)

// Update replace current values from domain.TModel by passed/*
func (repo *mongoRepo) Update(todo *domain.TModel) error {
	var item *domain.MongoModel
	utils.CopyStruct(todo, &item)

	ctx, cancel := utils.NewContextWithTimeout(timeoutIn)
	defer cancel()

	log.Infof("updating todo-list: %v", utils.ToString(item))
	log.Infof("id: %v", item.ID.Hex())
	res, err := repo.collection().ReplaceOne(ctx, bson.M{"_id": item.ID}, item)

	if err != nil || (res.MatchedCount == 0 && res.UpsertedCount == 0) {
		log.Errorf("error updating document: %v", err)
		return DocumentNotUpdated
	}

	return nil
}
