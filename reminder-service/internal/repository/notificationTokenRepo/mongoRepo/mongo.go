package mongoRepo

import (
	"context"
	"time"
	"tosinjs/reminder-service/internal/entity/errorEntity"
	"tosinjs/reminder-service/internal/entity/notificationTokenEntity"
	"tosinjs/reminder-service/internal/repository/notificationTokenRepo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoRepo struct {
	mColl *mongo.Collection
	ctx   context.Context
}

func New(mColl *mongo.Collection, ctx context.Context) notificationTokenRepo.NotificationTokenRepository {
	return mongoRepo{
		mColl: mColl,
		ctx:   ctx,
	}
}

func (m mongoRepo) GetNotificationTokens(userId string) ([]string, *errorEntity.LayerError) {
	filter := bson.D{{Key: "userId", Value: userId}}
	notifTokenCursor, err := m.mColl.Find(m.ctx, filter)
	if err != nil {
		return nil, errorEntity.InternalServerError("repo", err)
	}

	notifTokens := make([]string, 0)
	for notifTokenCursor.Next(m.ctx) {
		var notifToken notificationTokenEntity.NotificationToken
		err = notifTokenCursor.Decode(&notifToken)
		notifTokens = append(notifTokens, notifToken.DeviceId)
	}

	if err != nil {
		return nil, errorEntity.InternalServerError("repo", err)
	}
	return notifTokens, nil
}

func (m mongoRepo) RegisterNotificationId(notifDetails notificationTokenEntity.NotificationToken) *errorEntity.LayerError {
	filter := bson.D{{Key: "deviceId", Value: notifDetails.DeviceId}}
	res := m.mColl.FindOne(m.ctx, filter)

	if res.Err() != nil {
		if res.Err() == mongo.ErrNoDocuments {
			notifDetails.ID = primitive.NewObjectID()
			_, err := m.mColl.InsertOne(m.ctx, notifDetails)
			return errorEntity.InternalServerError("repo", err)
		}
		return errorEntity.InternalServerError("repo", res.Err())
	}
	_, err := m.mColl.UpdateOne(m.ctx, filter, bson.M{"$set": bson.M{"timestamp": time.Now().UTC()}})
	if err != nil {
		return errorEntity.InternalServerError("repo", err)
	}
	return nil
}
