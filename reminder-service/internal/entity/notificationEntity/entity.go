package notificationEntity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateNotificationReq struct {
	UserId       string `json:"userId" validate:"required"`
	Notification string `json:"notification" validate:"required"`
}

type Notification struct {
	ID           primitive.ObjectID `bson:"_id" json:"id"`
	UserId       string             `bson:"userId" json:"userId"`
	Notification string             `bson:"notification" json:"notification"`
	CreatedAt    time.Time          `bson:"createdAt" json:"createdAt"`
}
