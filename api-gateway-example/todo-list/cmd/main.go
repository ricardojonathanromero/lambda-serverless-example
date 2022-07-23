package main

import (
	ddLambda "github.com/DataDog/datadog-lambda-go"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/infrastructure/mongo"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-list/internal/handler"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-list/internal/repository"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-list/internal/service"
	log "github.com/sirupsen/logrus"
)

func main() {
	// init mongodb connection
	conn, err := mongo.NewConn()
	if err != nil {
		log.Fatalf("error initializing mongodb conn: %v", err)
	}

	hdl := handler.New(service.New(repository.New(conn)))

	// initialize lambda wrapper
	lambda.Start(ddLambda.WrapFunction(hdl, nil))
}
