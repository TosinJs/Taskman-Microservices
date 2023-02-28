package routes

import (
	"tosinjs/notification-service/cmd/api/handlers/notificationHandler"
	"tosinjs/notification-service/cmd/api/middleware/authMiddleware"
	"tosinjs/notification-service/internal/service/authService"
	"tosinjs/notification-service/internal/service/notificationService"
	"tosinjs/notification-service/internal/service/notificationTokenService"
	"tosinjs/notification-service/internal/service/validationService"

	"github.com/gin-gonic/gin"
)

func NotificationRoutes(
	v1 *gin.RouterGroup,
	notifSVC notificationService.NotificationService,
	notifTokenSVC notificationTokenService.NotificationTokenService,
	authSVC authService.AuthService,
	validationSVC validationService.ValidationService,
) {
	authMiddleware := authMiddleware.New(authSVC)
	notifHandler := notificationHandler.NewHandler(notifSVC, notifTokenSVC, validationSVC)
	notifRoutes := v1.Group("/notifications")

	notifRoutes.Use(authMiddleware.VerifyJWT())

	notifRoutes.GET("", notifHandler.GetNotifications)
	notifRoutes.POST("", notifHandler.CreateNotification)
	notifRoutes.POST("notification_token", notifHandler.RegisterNotificationId)
}
