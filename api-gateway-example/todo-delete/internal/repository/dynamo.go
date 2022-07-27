package repository

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-delete/internal/port"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/utils"
)

type dynamoRepo struct {
	client    *dynamodb.Client
	tableName string
}

var _ port.IRepository = (*dynamoRepo)(nil)

// Delete removes domain.TModel from db/*
func (repo *dynamoRepo) Delete(id interface{}) error {
	ctx, cancel := utils.NewContextWithTimeout(timeoutIn)
	defer cancel()

	log.Trace("deleting item...")
	out, err := repo.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		Key: map[string]types.AttributeValue{
			"Id": &types.AttributeValueMemberN{Value: fmt.Sprintf("%v", id.(int64))},
		},
		TableName: aws.String(repo.tableName),
	})

	if err != nil {
		log.Errorf("document not deleted: %v", err)
		return DocumentNotDeleted
	}

	log.Tracef("document %v deleted", id)
	log.Info(out)
	return nil
}
