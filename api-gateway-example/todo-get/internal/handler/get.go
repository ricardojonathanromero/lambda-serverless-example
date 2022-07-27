package handler

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-get/internal/domain"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-get/internal/port"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/utils"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"net/http"
)

var (
	log           = logrus.New()
	commonHeaders = map[string]string{"Content-Type": "application/json"}
)

type handle struct {
	srv port.IService
}

var _ port.IHandle = (*handle)(nil)

// HandleRequest business logic/*
func (h *handle) HandleRequest(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var id interface{}
	res := events.APIGatewayProxyResponse{Headers: commonHeaders, IsBase64Encoded: true}

	_, _ = tracer.StartSpanFromContext(ctx, "get todo")
	log.Info("start HandleRequest()")
	defer log.Info("end HandleRequest()")

	idStr := req.PathParameters["id"]
	log.Infof("id: %v", idStr)

	if !primitive.IsValidObjectID(idStr) && utils.StringToInt(idStr) <= 0 {
		log.Errorf("invalid id: %v", id)
		res.StatusCode = http.StatusBadRequest
		res.Body = utils.EncodeStr(utils.ToString(domain.NewErr("invalid_request", "the id sent is not valid")))
		return res, nil
	}

	if primitive.IsValidObjectID(idStr) {
		id, _ = primitive.ObjectIDFromHex(idStr)
	} else {
		id = utils.StringToInt(idStr)
	}

	log.Info("searching...")
	item, err := h.srv.GetTodo(id)
	if err != nil {
		log.Errorf("error retrieving item: %v", err)
		res.StatusCode = http.StatusNotFound
		res.Body = utils.EncodeStr(utils.ToString(domain.NewErr("not_found", err.Error())))
		return res, nil
	}

	log.Info("returning response")
	res.StatusCode = http.StatusOK
	res.Body = utils.EncodeStr(utils.ToString(item))
	return res, nil
}

//New constructor /*
func New(srv port.IService) port.IHandle {
	return &handle{srv: srv}
}
