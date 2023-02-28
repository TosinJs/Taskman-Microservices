package notificationRepo

import (
	"tosinjs/notification-service/internal/entity/errorEntity"
	"tosinjs/notification-service/internal/entity/notificationEntity"
)

type NotificationRepository interface {
	CreateNotification(notif notificationEntity.Notification) *errorEntity.LayerError
	GetNotifications(userId string) ([]notificationEntity.Notification, *errorEntity.LayerError)
}
