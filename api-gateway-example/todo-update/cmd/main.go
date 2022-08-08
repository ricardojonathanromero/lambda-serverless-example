package main

import (
	ddLambda "github.com/DataDog/datadog-lambda-go"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-update/internal/handler"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-update/internal/port"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-update/internal/repository"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-update/internal/service"
	"github.com/ricardojonathanromero/lambda-utilities/dbs/dynamodb"
	"github.com/ricardojonathanromero/lambda-utilities/dbs/mongo"
	"github.com/ricardojonathanromero/lambda-utilities/utils"
	log "github.com/sirupsen/logrus"
	"gopkg.in/DataDog/dd-trace-go.v1/profiler"
	"strings"
	"time"
)

func main() {
	// profiles
	err := profiler.Start(
		profiler.WithService("todo-update"),
		profiler.WithEnv(utils.GetEnv("ENV", "local")),
		profiler.WithVersion("v1.0.0"),
		profiler.WithTags("cloud: aws"),
		profiler.WithProfileTypes(
			profiler.CPUProfile,
			profiler.HeapProfile,
		),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer profiler.Stop()

	// configuring logger formatter
	if strings.EqualFold(utils.GetEnv("CLOUD", "false"), "true") {
		log.SetFormatter(&log.JSONFormatter{TimestampFormat: time.RFC3339})
	}

	// configuring repository
	var repo port.IRepository
	const db string = "todos"
	if strings.EqualFold(utils.GetEnv("DB", "mongodb"), "dynamodb") {
		conn, err := dynamodb.NewSess()
		if err != nil {
			log.Fatalf("error initializing dynamodb conn: %v", err)
		}
		repo = repository.NewDynamo(conn, db)
	} else {
		conn, err := mongo.NewConn()
		if err != nil {
			log.Fatalf("error initializing mongodb conn: %v", err)
		}
		repo = repository.NewMongo(conn, db, db)
	}

	hdl := handler.New(service.New(repo))

	// initialize lambda wrapper
	lambda.Start(ddLambda.WrapFunction(hdl.HandleRequest, nil))
}
