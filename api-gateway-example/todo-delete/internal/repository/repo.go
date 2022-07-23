package repository

import (
	"errors"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-delete/internal/port"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/utils"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

const timeoutIn = 10 * time.Second

var (
	log                = logrus.New()
	DocumentNotDeleted = errors.New("document could not be deleted")
)

type repository struct {
	conn *mongo.Client
	db   string
	col  string
}

var _ port.IRepository = (*repository)(nil)

// Delete removes domain.TModel from db/*
func (repo *repository) Delete(id primitive.ObjectID) error {
	ctx, cancel := utils.NewContextWithTimeout(timeoutIn)
	defer cancel()

	res, err := repo.collection().DeleteOne(ctx, bson.M{"_id": id})
	if err != nil || (res.DeletedCount == 0) {
		log.Errorf("document not deleted: %v", err)
		return DocumentNotDeleted
	}

	log.Tracef("document %v deleted", id)
	return nil
}

// collection returns mongodb collection/*
func (repo *repository) collection() *mongo.Collection {
	return repo.conn.Database(repo.db).Collection(repo.col)
}

// New constructor for repository package/*
func New(conn *mongo.Client) port.IRepository {
	return &repository{conn: conn}
}
