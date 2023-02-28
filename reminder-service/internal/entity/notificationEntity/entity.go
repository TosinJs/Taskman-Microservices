package notificationEntity

type CreateNotificationReq struct {
	UserId       string `json:"userId" validate:"required"`
	Notification string `json:"notification" validate:"required"`
}
