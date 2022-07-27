package repository

import (
	"errors"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-update/internal/port"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

const timeoutIn = 10 * time.Second

var (
	log                = logrus.New()
	DocumentNotUpdated = errors.New("document could not be updated")
)

// collection returns mongodb collection/*
func (repo *mongoRepo) collection() *mongo.Collection {
	return repo.conn.Database(repo.db).Collection(repo.col)
}

// NewMongo constructor for repository package/*
func NewMongo(conn *mongo.Client, db string, col string) port.IRepository {
	return &mongoRepo{conn: conn, db: db, col: col}
}

// NewDynamo constructor for repository package/*
func NewDynamo(client *dynamodb.Client, tableName string) port.IRepository {
	return &dynamoRepo{client: client, tableName: tableName}
}
