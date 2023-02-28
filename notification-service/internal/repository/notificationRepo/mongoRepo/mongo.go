package mongoRepo

import (
	"context"
	"tosinjs/notification-service/internal/entity/errorEntity"
	"tosinjs/notification-service/internal/entity/notificationEntity"
	"tosinjs/notification-service/internal/repository/notificationRepo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoRepo struct {
	mColl *mongo.Collection
	ctx   context.Context
}

func New(mColl *mongo.Collection, ctx context.Context) notificationRepo.NotificationRepository {
	return mongoRepo{
		mColl: mColl,
		ctx:   ctx,
	}
}

func (m mongoRepo) CreateNotification(notif notificationEntity.Notification) *errorEntity.LayerError {
	notif.ID = primitive.NewObjectID()
	_, err := m.mColl.InsertOne(m.ctx, notif)
	if err != nil {
		return errorEntity.InternalServerError("repo", err)
	}
	return nil
}

func (m mongoRepo) GetNotifications(userId string) ([]notificationEntity.Notification, *errorEntity.LayerError) {
	filter := bson.D{{Key: "userId", Value: userId}}
	notifCursor, err := m.mColl.Find(m.ctx, filter)

	if err != nil {
		return nil, errorEntity.InternalServerError("repo", err)
	}

	defer notifCursor.Close(m.ctx)

	notifications := make([]notificationEntity.Notification, 0)
	for notifCursor.Next(m.ctx) {
		var notification notificationEntity.Notification
		err = notifCursor.Decode(&notification)
		notifications = append(notifications, notification)
	}

	if err != nil {
		return nil, errorEntity.InternalServerError("repo", err)
	}

	return notifications, nil
}
