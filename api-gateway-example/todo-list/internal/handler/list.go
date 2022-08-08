package handler

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-list/internal/domain"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-list/internal/port"
	"github.com/ricardojonathanromero/lambda-utilities/utils"
	"github.com/sirupsen/logrus"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"net/http"
)

var (
	log           = logrus.New()
	commonHeaders = map[string]string{"Content-Type": "application/json"}
)

type handler struct {
	srv port.IService
}

var _ port.IHandler = (*handler)(nil)

// HandleRequest business logic/*
func (h *handler) HandleRequest(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	res := events.APIGatewayProxyResponse{Headers: commonHeaders, IsBase64Encoded: true}

	_, _ = tracer.StartSpanFromContext(ctx, "list todos")
	log.Info("start HandleRequest()")
	defer log.Info("end HandleRequest()")

	limit := utils.StringToInt(req.QueryStringParameters["limit"])
	offset := utils.StringToInt(req.QueryStringParameters["offset"])

	if limit <= 0 {
		limit = 10
	}

	log.Infof("limit: %v, offset: %v", limit, offset)
	result, err := h.srv.GetAll(limit, offset)
	if err != nil {
		res.StatusCode = http.StatusConflict
		res.Body = utils.EncodeStr(utils.ToString(domain.NewErr("not_found", err.Error())))
		log.Error("error result GetAll()")
		return res, nil
	}

	log.Info("returning response")
	res.StatusCode = http.StatusOK
	res.Body = utils.EncodeStr(utils.ToString(result))
	return res, nil
}

func New(srv port.IService) port.IHandler {
	return &handler{srv: srv}
}
