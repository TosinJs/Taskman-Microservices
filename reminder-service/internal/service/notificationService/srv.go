package notificationService

import (
	"context"
	"encoding/json"
	"fmt"
	"tosinjs/reminder-service/internal/entity/errorEntity"
	"tosinjs/reminder-service/internal/entity/notificationEntity"

	amqp "github.com/rabbitmq/amqp091-go"
)

type notificationService struct {
	ctx   context.Context
	nq    amqp.Queue
	nChan *amqp.Channel
}

type NotificationService interface {
	PublishNotification(notificationEntity.CreateNotificationReq) *errorEntity.LayerError
}

func New(
	ctx context.Context,
	nq amqp.Queue,
	nChan *amqp.Channel,
) NotificationService {
	return notificationService{
		ctx:   ctx,
		nq:    nq,
		nChan: nChan,
	}
}

func (n notificationService) PublishNotification(
	notification notificationEntity.CreateNotificationReq,
) *errorEntity.LayerError {
	body, err := json.Marshal(notification)
	if err != nil {
		return errorEntity.InternalServerError("service", err)
	}
	err = n.nChan.PublishWithContext(
		n.ctx,
		"",
		n.nq.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		},
	)
	if err != nil {
		fmt.Println(err)
		return errorEntity.InternalServerError("service", err)
	}
	fmt.Println("published")
	return nil
}
