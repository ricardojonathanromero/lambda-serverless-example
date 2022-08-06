package repository

import (
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-create/internal/domain"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-create/internal/port"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoRepo struct {
	conn *mongo.Client
	db   string
	col  string
}

var _ port.IRepository = (*mongoRepo)(nil)

// Insert inserts a new domain.TModel document into db/*
func (repo *mongoRepo) Insert(todo domain.TModel) (string, error) {
	var item *domain.MongoModel
	utils.CopyStruct(todo, &item)

	ctx, cancel := utils.NewContextWithTimeout(timeoutIn)
	defer cancel()

	log.Tracef("inserting todo-list object: %v", utils.ToString(todo))
	inserted, err := repo.collection().InsertOne(ctx, item)
	if err != nil {
		return primitive.NilObjectID.Hex(), DocumentNotInserted
	}

	log.Tracef("id inserted: %v", inserted.InsertedID.(primitive.ObjectID))
	return inserted.InsertedID.(primitive.ObjectID).Hex(), nil
}
