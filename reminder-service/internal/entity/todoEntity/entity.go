package todoEntity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateTodoReq struct {
	UserId    string `json:"userId" validate:"required"`
	Todo      string `json:"todo" validate:"required"`
	CreatedAt string `json:"createdAt"`
	Due       string `json:"due" validate:"required"`
	RemindMe  string `json:"remindMe"`
	Recurring bool   `json:"recurring" validate:"required"`
}

type Todo struct {
	Id        primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	UserId    string             `bson:"userId" json:"userId,omitempty"`
	Todo      string             `bson:"todo" json:"todo,omitempty"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt,omitempty"`
	Due       time.Time          `bson:"due" json:"due,omitempty"`
	RemindMe  time.Time          `bson:"remindMe" json:"remindMe,omitempty"`
	Recurring bool               `bson:"recurring" json:"recurring"`
	Done      bool               `bson:"done" json:"done"`
}
