package notificationService

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"tosinjs/notification-service/internal/entity/errorEntity"
	"tosinjs/notification-service/internal/entity/notificationEntity"
	"tosinjs/notification-service/internal/repository/notificationRepo"
	"tosinjs/notification-service/internal/service/notificationTokenService"

	"firebase.google.com/go/v4/messaging"
	"github.com/rabbitmq/amqp091-go"
	amqp "github.com/rabbitmq/amqp091-go"
)

type notificationService struct {
	notifRepo     notificationRepo.NotificationRepository
	notifTokenSVC notificationTokenService.NotificationTokenService
	fcm           *messaging.Client
	ctx           context.Context
	nq            amqp.Queue
	nChan         *amqp.Channel
}

type NotificationService interface {
	CreateNotification(notificationEntity.CreateNotificationReq) *errorEntity.LayerError
	GetNotifications(userId string) ([]notificationEntity.Notification, *errorEntity.LayerError)
	SendNotification(notifToken, message string) *errorEntity.LayerError
	SendBatchNotifications(notifTokens []string, message string) *errorEntity.LayerError
	ConsumeNotification(nmsg amqp.Delivery) *errorEntity.LayerError
}

func New(
	notifRepo notificationRepo.NotificationRepository,
	notifTokenSVC notificationTokenService.NotificationTokenService,
	fcm *messaging.Client,
	ctx context.Context,
	nq amqp.Queue,
	nChan *amqp.Channel,
) NotificationService {
	return notificationService{
		notifRepo:     notifRepo,
		notifTokenSVC: notifTokenSVC,
		fcm:           fcm,
		ctx:           ctx,
		nq:            nq,
		nChan:         nChan,
	}
}

func (n notificationService) ConsumeNotification(nmsg amqp091.Delivery) *errorEntity.LayerError {

	var notification notificationEntity.CreateNotificationReq
	err := json.Unmarshal(nmsg.Body, &notification)
	if err != nil {
		return errorEntity.InternalServerError("service", err)
	}
	svcErr := n.CreateNotification(notification)
	if svcErr != nil {
		return svcErr
	}
	tokens, svcErr := n.notifTokenSVC.GetNotificationTokens(notification.UserId)
	if svcErr != nil {
		return svcErr
	}
	svcErr = n.SendBatchNotifications(tokens, notification.Notification)
	if svcErr != nil {
		fmt.Println(err)
	}
	return nil
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
