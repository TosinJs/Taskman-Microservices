package notificationTokenEntity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RegisterNotifIdReq struct {
	UserId   string `json:"userId" validate:"required"`
	DeviceId string `json:"deviceId" validate:"required"`
}

type NotificationToken struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	UserId    string             `bson:"userId" json:"userId"`
	DeviceId  string             `bson:"deviceId" json:"deviceId"`
	Timestamp time.Time          `bson:"timestamp" json:"timestamp"`
}
