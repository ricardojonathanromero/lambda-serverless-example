package handler

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/go-playground/validator/v10"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-create/internal/domain"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-create/internal/port"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/utils"
	"github.com/sirupsen/logrus"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"net/http"
)

type handle struct {
	srv port.IService
}

var (
	log                         = logrus.New()
	commonHeaders               = map[string]string{"Content-Type": "application/json"}
	_             port.IHandler = (*handle)(nil)
)

// HandleRequest handle request/*
func (h *handle) HandleRequest(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var payload *domain.Req
	res := events.APIGatewayProxyResponse{Headers: commonHeaders, IsBase64Encoded: true}

	_, _ = tracer.StartSpanFromContext(ctx, "list todos")
	log.Info("start HandleRequest()")
	defer log.Info("end HandleRequest()")

	// get body request
	utils.StringToStruct(req.Body, &payload)

	log.Info("validating payload...")
	if err := validator.New().Struct(payload); err != nil {
		log.Errorf("error request: %v", err)
		res.StatusCode = http.StatusBadRequest
		res.Body = utils.EncodeStr(utils.ToString(domain.NewErr("invalid_request", "payload is not valid")))
		return res, nil
	}

	log.Info("removing...")
	id, err := h.srv.CreateTodo(payload.Name, payload.Priority)
	if err != nil {
		log.Errorf("error retrieving item: %v", err)
		res.StatusCode = http.StatusNotAcceptable
		res.Body = utils.EncodeStr(utils.ToString(domain.NewErr("task_not_created", err.Error())))
		return res, nil
	}

	log.Info("returning response")
	res.StatusCode = http.StatusCreated
	res.Body = utils.EncodeStr(utils.ToString(&domain.Res{ID: id}))
	return res, nil
}

// New constructor/*
func New(srv port.IService) port.IHandler {
	return &handle{srv: srv}
}
