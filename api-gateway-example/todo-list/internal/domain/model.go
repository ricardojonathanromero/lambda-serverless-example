package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Result struct {
	Data     []TModel `json:"data"`
	Metadata Metadata `json:"metadata"`
}

type Metadata struct {
	Limit  int64
	Offset int64
	Total  int64
}

// TModel describes the object /*
type TModel struct {
	// ID the identifier generated by db
	ID primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	// Name is the name of the card
	Name string `json:"name,omitempty" bson:"name,omitempty"`
	// Priority must be one of these: ['urgent', 'high', 'normal', 'medium', 'low']
	Priority string `json:"priority,omitempty" bson:"priority,omitempty"`
	// CreatedAt describes when the card has been created
	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	// UpdatedAt describes when the card has been updated
	UpdatedAt time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type ErrRes struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func NewErr(code, msg string) *ErrRes {
	return &ErrRes{Code: code, Message: msg}
}