package notificationTokenRepo

import (
	"tosinjs/notification-service/internal/entity/errorEntity"
	"tosinjs/notification-service/internal/entity/notificationTokenEntity"
)

type NotificationTokenRepository interface {
	GetNotificationTokens(userId string) ([]string, *errorEntity.LayerError)
	RegisterNotificationId(notifDetails notificationTokenEntity.NotificationToken) *errorEntity.LayerError
}
