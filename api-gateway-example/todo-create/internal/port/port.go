package port

import (
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/todo-create/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IRepository interface {
	Inset(todo domain.TModel) (primitive.ObjectID, error)
}
