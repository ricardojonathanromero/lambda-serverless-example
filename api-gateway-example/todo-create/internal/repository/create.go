package repository

import (
	"errors"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-create/internal/domain"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-create/internal/port"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/utils"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

const timeoutIn = 10 * time.Second

var (
	log                 = logrus.New()
	DocumentNotInserted = errors.New("document not inserted")
)

type repository struct {
	conn *mongo.Client
	db   string
	col  string
}

var _ port.IRepository = (*repository)(nil)

// collection returns mongodb collection/*
func (repo *repository) collection() *mongo.Collection {
	return repo.conn.Database(repo.db).Collection(repo.col)
}

// Inset inserts a new domain.TModel document into db/*
func (repo *repository) Inset(todo domain.TModel) (primitive.ObjectID, error) {
	ctx, cancel := utils.NewContextWithTimeout(timeoutIn)
	defer cancel()

	log.Tracef("inserting todo-list object: %v", utils.ToString(todo))
	inserted, err := repo.collection().InsertOne(ctx, todo)
	if err != nil {
		return primitive.NilObjectID, DocumentNotInserted
	}

	log.Tracef("id inserted: %v", inserted.InsertedID.(primitive.ObjectID))
	return inserted.InsertedID.(primitive.ObjectID), nil
}

// New constructor for repository package/*
func New(conn *mongo.Client) port.IRepository {
	return &repository{conn: conn}
}
