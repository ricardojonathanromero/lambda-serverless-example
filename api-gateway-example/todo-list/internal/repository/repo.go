package repository

import (
	"errors"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-list/internal/domain"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-list/internal/port"
	utils2 "github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/utils"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const timeoutIn = 10 * time.Second

var (
	log               = logrus.New()
	DocumentsNotFound = errors.New("documents not found")
)

type repository struct {
	conn *mongo.Client
	db   string
	col  string
}

var _ port.IRepository = (*repository)(nil)

// CountDocuments returns total documents/*
func (repo *repository) CountDocuments() (int64, error) {
	ctx, cancel := utils2.NewContextWithTimeout(timeoutIn)
	defer cancel()

	return repo.collection().CountDocuments(ctx, bson.M{})
}

// FindAll returns all documents from db/*
func (repo *repository) FindAll(limit, offset int64) ([]domain.TModel, error) {
	var res []domain.TModel

	log.Trace("init context")
	ctx, cancel := utils2.NewContextWithTimeout(timeoutIn)
	defer cancel()

	log.Trace("finding in db")
	opts := options.Find().SetLimit(limit).SetSkip(offset).SetSort("_id")

	result, err := repo.collection().Find(ctx, bson.M{}, opts)
	if err != nil {
		log.Errorf("documents not found: %v", err)
		return res, DocumentsNotFound
	}

	err = result.All(ctx, &res)
	if err != nil {
		log.Errorf("no documents for decode: %v", err)
		return res, DocumentsNotFound
	}

	log.Tracef("document founded: %v", len(res))
	if len(res) <= 0 {
		res = make([]domain.TModel, 0)
		return res, nil
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
