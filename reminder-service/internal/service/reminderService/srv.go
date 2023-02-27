package reminderService

import (
	"fmt"
	"time"
	"tosinjs/reminder-service/internal/entity/notificationEntity"
	"tosinjs/reminder-service/internal/service/notificationService"
	"tosinjs/reminder-service/internal/service/notificationTokenService"

	"github.com/go-co-op/gocron"
)

type reminderService struct {
	sch           *gocron.Scheduler
	notifSVC      notificationService.NotificationService
	notifTokenSVC notificationTokenService.NotificationTokenService
}

type ReminderService interface {
	CreateReminder(userId, todoId, message string, t time.Time)
	CreateRecurringReminder(repeat int, unit string)
	DeleteReminder(todoId string)
}

func New(
	sch *gocron.Scheduler,
	notifSVC notificationService.NotificationService,
	notifTokenSVC notificationTokenService.NotificationTokenService,
) ReminderService {
	return reminderService{
		sch:           sch,
		notifSVC:      notifSVC,
		notifTokenSVC: notifTokenSVC,
	}
}

func (r reminderService) CreateReminder(userId, todoId, message string, time time.Time) {
	reminderFunc := func() {
		fmt.Println("got heres")
		svcErr := r.notifSVC.CreateNotification(notificationEntity.CreateNotificationReq{
			UserId:       userId,
			Notification: message,
		})
		if svcErr != nil {
			fmt.Println(svcErr)
		}
		notifTokens, svcErr := r.notifTokenSVC.GetNotificationTokens(userId)
		if svcErr != nil {
			fmt.Println(svcErr)
		}
		svcErr = r.notifSVC.SendBatchNotifications(notifTokens, message)
		if svcErr != nil {
			fmt.Println(svcErr)
		}
	}
	job, err := r.sch.Every(1).Minutes().Tag(todoId).StartAt(time).Do(reminderFunc)
	if err != nil {
		fmt.Println(err)
	}
	job.LimitRunsTo(1)
}

func (r reminderService) DeleteReminder(todoId string) {
	r.sch.RemoveByTag(todoId)
	return
}

func (r reminderService) CreateRecurringReminder(repeat int, unit string) {
	reminderFunc := func() {

	}
	switch unit {
	case "Day":
		r.sch.Every(repeat).Day().Do(reminderFunc)
	case "Week":
		r.sch.Every(repeat).Week().Do(reminderFunc)
	case "Month":
		r.sch.Every(repeat).Month().Do(reminderFunc)
	default:
		r.sch.Every(repeat).Monday().Do(reminderFunc)
	}
}
