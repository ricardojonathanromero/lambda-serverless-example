package repository

import (
	"errors"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-list/internal/port"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

const (
	timeoutIn = 5 * time.Second
	idIndex   = "Id-Index"
)

var (
	log               = logrus.New()
	DocumentsNotFound = errors.New("documents not found")
)

// NewDynamoRepo constructor for repo kind dynamodb format/*
func NewDynamoRepo(client *dynamodb.Client, tableName string) port.IRepository {
	return &dynamoRepository{client: client, tableName: tableName}
}

// collection returns mongodb collection/*
func (repo *mongoRepository) collection() *mongo.Collection {
	return repo.conn.Database(repo.db).Collection(repo.col)
}

// NewMongoRepo constructor for repository package/*
func NewMongoRepo(conn *mongo.Client, db, col string) port.IRepository {
	return &mongoRepository{conn: conn, db: db, col: col}
}
