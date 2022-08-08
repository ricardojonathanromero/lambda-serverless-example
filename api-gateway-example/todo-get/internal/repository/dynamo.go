package repository

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-get/internal/domain"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-get/internal/port"
	"github.com/ricardojonathanromero/lambda-utilities/utils"
)

type dynamoRepo struct {
	conn      *dynamodb.Client
	tableName string
}

var _ port.IRepository = (*dynamoRepo)(nil)

// FindById returns document type domain.TModel from db/*
func (repo *dynamoRepo) FindById(id interface{}) (*domain.TModel, error) {
	var res *domain.TModel
	var item *domain.DynamoModel

	log.Trace("looking for object")
	// init context.Context
	ctx, cancel := utils.NewContextWithTimeout(timeoutIn)
	defer cancel()

	log.Trace("query")
	out, err := repo.conn.GetItem(ctx, &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"Id": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", id.(int64))},
		},
		TableName: aws.String(repo.tableName),
	})

	if err != nil || out.Item == nil {
		log.Errorf("error retrieving object: %v%v", err, out.Item)
		return res, DocumentNotFound
	}

	utils.DynamoMapToInterface(out.Item, &item)

	res = &domain.TModel{
		ID:        item.ID,
		Name:      item.Name,
		Status:    item.Status,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}

	log.Tracef("result: %v", res)
	return res, nil
}
