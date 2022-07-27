package repository

import (
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-delete/internal/port"
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

// Delete removes domain.TModel from db/*
func (repo *mongoRepo) Delete(id interface{}) error {
	ctx, cancel := utils.NewContextWithTimeout(timeoutIn)
	defer cancel()

	log.Trace("deleting item...")
	res, err := repo.collection().DeleteOne(ctx, bson.M{"_id": id.(primitive.ObjectID)})
	if err != nil || (res.DeletedCount == 0) {
		log.Errorf("document not deleted: %v%v", err, res)
		return DocumentNotDeleted
	}

	log.Tracef("document %v deleted", id)
	return nil
}
