package notificationTokenRepo

import (
	"tosinjs/reminder-service/internal/entity/errorEntity"
	"tosinjs/reminder-service/internal/entity/notificationTokenEntity"
)

type NotificationTokenRepository interface {
	GetNotificationTokens(userId string) ([]string, *errorEntity.LayerError)
	RegisterNotificationId(notifDetails notificationTokenEntity.NotificationToken) *errorEntity.LayerError
}
