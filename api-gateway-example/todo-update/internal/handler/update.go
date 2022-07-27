package handler

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-update/internal/domain"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-update/internal/port"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/utils"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	var payload *domain.TModel
	res := events.APIGatewayProxyResponse{Headers: commonHeaders, IsBase64Encoded: true}

	_, _ = tracer.StartSpanFromContext(ctx, "update todo")
	log.Info("start HandleRequest()")
	defer log.Info("end HandleRequest()")

	utils.StringToStruct(req.Body, &payload)

	log.Info("validating payload...")
	if payload == nil {
		log.Errorf("error request: %v", payload)
		res.StatusCode = http.StatusBadRequest
		res.Body = utils.EncodeStr(utils.ToString(domain.NewErr("invalid_request", "payload is not valid")))
		return res, nil
	}

	idStr := fmt.Sprintf("%v", payload.ID)
	if !primitive.IsValidObjectID(idStr) && utils.StringToInt(idStr) <= 0 {
		log.Errorf("invalid id: %v", idStr)
		res.StatusCode = http.StatusBadRequest
		res.Body = utils.EncodeStr(utils.ToString(domain.NewErr("invalid_request", "the id sent is not valid")))
		return res, nil
	}

	log.Info("updating item...")
	err := h.srv.UpdateItem(payload)
	if err != nil {
		log.Errorf("error updating item: %v", err)
		res.StatusCode = http.StatusConflict
		res.Body = utils.EncodeStr(utils.ToString(domain.NewErr("not_updated", err.Error())))
		return res, nil
	}

	log.Info("item updated")
	res.StatusCode = http.StatusOK
	res.Body = utils.EncodeStr(`{"message":"updated"}`)
	return res, nil
}

// New constructor/*
func New(srv port.IService) port.IHandler {
	return &handle{srv: srv}
}
