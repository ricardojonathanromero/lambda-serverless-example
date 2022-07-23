package main

import (
	ddLambda "github.com/DataDog/datadog-lambda-go"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/infrastructure/mongo"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-list/internal/handler"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-list/internal/repository"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-list/internal/service"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/utils"
	log "github.com/sirupsen/logrus"
	"gopkg.in/DataDog/dd-trace-go.v1/profiler"
)

func main() {
	// profiles
	err := profiler.Start(
		profiler.WithService("todo-list"),
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

	// init mongodb connection
	conn, err := mongo.NewConn()
	if err != nil {
		log.Fatalf("error initializing mongodb conn: %v", err)
	}

	hdl := handler.New(service.New(repository.New(conn)))

	// initialize lambda wrapper
	lambda.Start(ddLambda.WrapFunction(hdl, nil))
}