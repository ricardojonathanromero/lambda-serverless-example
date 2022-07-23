package repository

import (
	"errors"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-get/internal/domain"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-get/internal/port"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/utils"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

const timeoutIn = 10 * time.Second

var (
	log              = logrus.New()
	DocumentNotFound = errors.New("document not found")
)

type repository struct {
	conn *mongo.Client
	db   string
	col  string
}

var _ port.IRepository = (*repository)(nil)

// FindById returns document type domain.TModel from db/*
func (repo *repository) FindById(id primitive.ObjectID) (domain.TModel, error) {
	var res domain.TModel

	ctx, cancel := utils.NewContextWithTimeout(timeoutIn)
	defer cancel()

	err := repo.collection().FindOne(ctx, bson.M{"_id": id}).Decode(&res)
	if err != nil {
		log.Errorf("document %v not found: %v", id, err)
		return res, DocumentNotFound
	}

	return res, nil
}

// collection returns mongodb collection/*
func (repo *repository) collection() *mongo.Collection {
	return repo.conn.Database(repo.db).Collection(repo.col)
}

// New constructor for repository package/*
func New(conn *mongo.Client) port.IRepository {
	return &repository{conn: conn}
}
