package repository

import (
	"errors"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-update/internal/domain"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-update/internal/port"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/utils"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

const timeoutIn = 10 * time.Second

var (
	log                = logrus.New()
	DocumentNotUpdated = errors.New("document could not be updated")
)

type repository struct {
	conn *mongo.Client
	db   string
	col  string
}

var _ port.IRepository = (*repository)(nil)

// Update replace current values from domain.TModel by passed/*
func (repo *repository) Update(todo domain.TModel) error {
	ctx, cancel := utils.NewContextWithTimeout(timeoutIn)
	defer cancel()

	log.Tracef("updating todo-list: %v", utils.ToString(todo))
	res, err := repo.collection().ReplaceOne(ctx, bson.M{"_id": todo.ID}, todo)
	if err != nil || (res.MatchedCount == 0 && res.UpsertedCount == 0) {
		log.Errorf("error updating document: %v", err)
		return DocumentNotUpdated
	}

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
