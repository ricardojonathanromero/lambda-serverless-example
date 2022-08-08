package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/ricardojonathanromero/lambda-serverless-example/sqs-example/event/internal/handler"
	"github.com/ricardojonathanromero/lambda-serverless-example/sqs-example/event/internal/service"
	"github.com/ricardojonathanromero/lambda-serverless-example/sqs-example/utils"
	log "github.com/sirupsen/logrus"
	"gopkg.in/DataDog/dd-trace-go.v1/profiler"
	"strings"
	"time"
)

func main() {
	// profiles
	err := profiler.Start(
		profiler.WithService("sqs-example"),
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

	hdl := handler.New(service.New())

	lambda.Start(hdl.HandleEvent)
}
