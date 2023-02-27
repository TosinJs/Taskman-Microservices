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
	UserId    string             `json:"userId,omitempty"`
	Todo      string             `json:"todo,omitempty"`
	CreatedAt time.Time          `json:"createdAt,omitempty"`
	Due       time.Time          `json:"due,omitempty"`
	RemindMe  time.Time          `json:"remindMe,omitempty"`
	Recurring bool               `json:"recurring"`
	Done      bool               `json:"done"`
}
