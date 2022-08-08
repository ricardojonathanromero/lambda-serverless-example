package repository

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-update/internal/domain"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-update/internal/port"
	"github.com/ricardojonathanromero/lambda-utilities/utils"
)

type dynamoRepo struct {
	client    *dynamodb.Client
	tableName string
}

var _ port.IRepository = (*dynamoRepo)(nil)

// Update replace current values from domain.TModel by passed/*
func (repo *dynamoRepo) Update(todo *domain.TModel) error {
	var item *domain.DynamoModel
	utils.CopyStruct(todo, &item)

	ctx, cancel := utils.NewContextWithTimeout(timeoutIn)
	defer cancel()

	log.Tracef("updating todo-list: %v", utils.ToString(todo))
	upd := &dynamodb.PutItemInput{Item: utils.ToDynamoDBMap(item), TableName: aws.String(repo.tableName)}
	_, err := repo.client.PutItem(ctx, upd)
	if err != nil {
		log.Errorf("error updating item: %v", err)
		return DocumentNotUpdated
	}

	log.Trace("document updated")
	return nil
}
