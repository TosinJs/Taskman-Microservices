package notificationTokenService

import (
	"time"
	"tosinjs/reminder-service/internal/entity/errorEntity"
	"tosinjs/reminder-service/internal/entity/notificationTokenEntity"
	"tosinjs/reminder-service/internal/repository/notificationTokenRepo"
)

type notificationTokenService struct {
	notifTokenRepo notificationTokenRepo.NotificationTokenRepository
}

type NotificationTokenService interface {
	GetNotificationTokens(userId string) ([]string, *errorEntity.LayerError)
	RegisterNotificationId(notifDetails notificationTokenEntity.RegisterNotifIdReq) *errorEntity.LayerError
}

func New(
	notifTokenRepo notificationTokenRepo.NotificationTokenRepository,
) NotificationTokenService {
	return notificationTokenService{
		notifTokenRepo: notifTokenRepo,
	}
}

func (n notificationTokenService) GetNotificationTokens(userId string) ([]string, *errorEntity.LayerError) {
	return n.notifTokenRepo.GetNotificationTokens(userId)
}

func (n notificationTokenService) RegisterNotificationId(
	notifDetails notificationTokenEntity.RegisterNotifIdReq,
) *errorEntity.LayerError {
	notificationToken := notificationTokenEntity.NotificationToken{
		UserId:    notifDetails.UserId,
		DeviceId:  notifDetails.DeviceId,
		Timestamp: time.Now().UTC(),
	}
	return n.notifTokenRepo.RegisterNotificationId(notificationToken)
}
