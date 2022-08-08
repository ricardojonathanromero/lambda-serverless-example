package repository

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-list/internal/domain"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-list/internal/port"
	"github.com/ricardojonathanromero/lambda-utilities/utils"
	"time"
)

type dynamoRepository struct {
	client    *dynamodb.Client
	tableName string
}

var _ port.IRepository = (*dynamoRepository)(nil)

// CountDocuments counts total of documents in DynamoDB/*
func (repo *dynamoRepository) CountDocuments() (int64, error) {
	ctx, cancel := utils.NewContextWithTimeout(timeoutIn)
	defer cancel()

	// count documents
	query := &dynamodb.ScanInput{Select: types.SelectCount, TableName: aws.String(repo.tableName)}
	out, err := repo.client.Scan(ctx, query)
	if err != nil {
		log.Errorf("error count documents: %v", err)
		return 0, err
	}

	return int64(out.Count), nil
}

// FindAll returns all documents paginated/*
func (repo *dynamoRepository) FindAll(limit int64, offset int64) ([]*domain.TModel, error) {
	var res []*domain.TModel
	var items []*domain.DynamoModel

	query := &dynamodb.ScanInput{
		TableName: aws.String(repo.tableName),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":createdAt": &types.AttributeValueMemberS{Value: time.Now().Format(time.RFC3339)},
		},
		FilterExpression: aws.String("CreatedAt < :createdAt"),
		IndexName:        aws.String(idIndex),
		Limit:            aws.Int32(int32(int(limit))),
		Select:           types.SelectAllAttributes,
	}

	if offset > 0 {
		query.ExclusiveStartKey = map[string]types.AttributeValue{"Id": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", offset)}}
	}

	pag := dynamodb.NewScanPaginator(repo.client, query)
	ctx, cancel := utils.NewContextWithTimeout(timeoutIn)
	defer cancel()

	if pag.HasMorePages() {
		out, err := pag.NextPage(ctx)
		if err != nil {
			log.Errorf("error next pages: %v", err)
			return res, DocumentsNotFound
		} else if len(out.Items) > 0 {
			utils.DynamoListToInterface(out.Items, &items)
			utils.CopyStruct(items, &res)
		} else {
			res = make([]*domain.TModel, 0)
		}
	}

	return res, nil
}
