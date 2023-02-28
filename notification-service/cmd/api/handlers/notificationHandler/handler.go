package notificationHandler

import (
	"net/http"
	"tosinjs/notification-service/internal/entity/notificationEntity"
	"tosinjs/notification-service/internal/entity/notificationTokenEntity"
	"tosinjs/notification-service/internal/entity/responseEntity"
	"tosinjs/notification-service/internal/service/notificationService"
	"tosinjs/notification-service/internal/service/notificationTokenService"
	"tosinjs/notification-service/internal/service/validationService"

	"github.com/gin-gonic/gin"
)

type notificationHandler struct {
	notifSVC      notificationService.NotificationService
	notifTokenSVC notificationTokenService.NotificationTokenService
	validationSVC validationService.ValidationService
}

func NewHandler(
	notifSVC notificationService.NotificationService,
	notifTokenSVC notificationTokenService.NotificationTokenService,
	validationSVC validationService.ValidationService,
) notificationHandler {
	return notificationHandler{
		notifSVC:      notifSVC,
		notifTokenSVC: notifTokenSVC,
		validationSVC: validationSVC,
	}
}

func (n notificationHandler) CreateNotification(c *gin.Context) {
	userId := c.GetString("userId")
	var req notificationEntity.CreateNotificationReq
	req.UserId = userId

	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, responseEntity.BuildErrorResponseObject(
			err.Error(), c.FullPath(),
		))
		return
	}
	if svcErr := n.validationSVC.Validate(req); svcErr != nil {
		c.AbortWithStatusJSON(svcErr.StatusCode, responseEntity.BuildServiceErrorResponseObject(
			svcErr, c.FullPath(),
		))
		return
	}
	if svcErr := n.notifSVC.CreateNotification(req); svcErr != nil {
		c.AbortWithStatusJSON(svcErr.StatusCode, responseEntity.BuildServiceErrorResponseObject(
			svcErr, c.FullPath(),
		))
		return
	}
	c.JSON(http.StatusAccepted, "success")
}

func (n notificationHandler) GetNotifications(c *gin.Context) {
	userId := c.GetString("userId")

	notifications, svcErr := n.notifSVC.GetNotifications(userId)
	if svcErr != nil {
		c.AbortWithStatusJSON(svcErr.StatusCode, responseEntity.BuildServiceErrorResponseObject(
			svcErr, c.FullPath(),
		))
		return
	}

	c.JSON(http.StatusAccepted, notifications)
}

func (n notificationHandler) RegisterNotificationId(c *gin.Context) {
	userId := c.GetString("userId")
	var req notificationTokenEntity.RegisterNotifIdReq
	req.UserId = userId

	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, responseEntity.BuildErrorResponseObject(
			err.Error(), c.FullPath(),
		))
		return
	}
	if svcErr := n.validationSVC.Validate(req); svcErr != nil {
		c.AbortWithStatusJSON(svcErr.StatusCode, responseEntity.BuildServiceErrorResponseObject(
			svcErr, c.FullPath(),
		))
		return
	}
	if svcErr := n.notifTokenSVC.RegisterNotificationId(req); svcErr != nil {
		c.AbortWithStatusJSON(svcErr.StatusCode, responseEntity.BuildServiceErrorResponseObject(
			svcErr, c.FullPath(),
		))
		return
	}
	c.JSON(http.StatusAccepted, "success")
}
