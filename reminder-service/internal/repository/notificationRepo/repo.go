package notificationRepo

import (
	"tosinjs/reminder-service/internal/entity/errorEntity"
	"tosinjs/reminder-service/internal/entity/notificationEntity"
)

type NotificationRepository interface {
	CreateNotification(notif notificationEntity.Notification) *errorEntity.LayerError
	GetNotifications(userId string) ([]notificationEntity.Notification, *errorEntity.LayerError)
}
