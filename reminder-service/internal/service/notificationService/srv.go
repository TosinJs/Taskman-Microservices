package notificationService

import (
	"context"
	"fmt"
	"time"
	"tosinjs/reminder-service/internal/entity/errorEntity"
	"tosinjs/reminder-service/internal/entity/notificationEntity"
	"tosinjs/reminder-service/internal/repository/notificationRepo"

	"firebase.google.com/go/v4/messaging"
)

type notificationService struct {
	notifRepo notificationRepo.NotificationRepository
	fcm       *messaging.Client
	ctx       context.Context
}

type NotificationService interface {
	CreateNotification(notificationEntity.CreateNotificationReq) *errorEntity.LayerError
	GetNotifications(userId string) ([]notificationEntity.Notification, *errorEntity.LayerError)
	SendNotification(notifToken, message string) *errorEntity.LayerError
	SendBatchNotifications(notifTokens []string, message string) *errorEntity.LayerError
}

func New(notifRepo notificationRepo.NotificationRepository, fcm *messaging.Client, ctx context.Context) NotificationService {
	return notificationService{
		notifRepo: notifRepo,
		fcm:       fcm,
		ctx:       ctx,
	}
}

func (n notificationService) CreateNotification(notif notificationEntity.CreateNotificationReq) *errorEntity.LayerError {
	notification := notificationEntity.Notification{
		UserId:       notif.UserId,
		Notification: notif.Notification,
		CreatedAt:    time.Now().UTC(),
	}
	return n.notifRepo.CreateNotification(notification)
}

func (n notificationService) GetNotifications(userId string) ([]notificationEntity.Notification, *errorEntity.LayerError) {
	return n.notifRepo.GetNotifications(userId)
}

func (n notificationService) SendNotification(notifToken, message string) *errorEntity.LayerError {
	res, err := n.fcm.Send(n.ctx, &messaging.Message{
		Token: notifToken,
		Data: map[string]string{
			message: message,
		},
	})

	if err != nil {
		return errorEntity.InternalServerError("service", err)
	}
	fmt.Println("Notification Sent:", res)
	return nil
}

func (n notificationService) SendBatchNotifications(notifTokens []string, message string) *errorEntity.LayerError {
	res, err := n.fcm.SendMulticast(n.ctx, &messaging.MulticastMessage{
		Data: map[string]string{
			message: message,
		},
		Tokens: notifTokens,
	})

	if err != nil {
		return errorEntity.InternalServerError("service", err)
	}
	fmt.Println("Notification Sent:", res)
	return nil
}
