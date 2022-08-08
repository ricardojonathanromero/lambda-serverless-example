package repository

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-create/internal/domain"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-create/internal/port"
	"github.com/ricardojonathanromero/lambda-utilities/utils"
)

type dynamoRepo struct {
	client    *dynamodb.Client
	tableName string
}

var _ port.IRepository = (*dynamoRepo)(nil)

// Insert inserts a new domain.TModel document into db/*
func (repo *dynamoRepo) Insert(todo domain.TModel) (string, error) {
	var item *domain.DynamoModel
	var result string
	utils.CopyStruct(todo, &item)

	ctx, cancel := utils.NewContextWithTimeout(timeoutIn)
	defer cancel()

	log.Tracef("inserting todo: %s", utils.ToString(todo))
	insert := &dynamodb.PutItemInput{Item: utils.ToDynamoDBMap(item), TableName: aws.String(repo.tableName)}
	res, err := repo.client.PutItem(ctx, insert)
	if err != nil {
		log.Errorf("error inserting item: %v", err)
		return result, DocumentNotInserted
	}

	log.Trace("document created")
	log.Info(utils.ToString(res.ResultMetadata))
	log.Info(utils.ToString(res.Attributes))
	return fmt.Sprintf("%d", todo.ID), nil
}
