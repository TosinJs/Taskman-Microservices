package routes

import (
	"tosinjs/reminder-service/cmd/api/handlers/todoHandler"
	"tosinjs/reminder-service/cmd/api/middleware/authMiddleware"
	"tosinjs/reminder-service/internal/service/authService"
	"tosinjs/reminder-service/internal/service/todoService"
	"tosinjs/reminder-service/internal/service/validationService"

	"github.com/gin-gonic/gin"
)

func TodoRoutes(
	v1 *gin.RouterGroup,
	todoSVC todoService.TodoService,
	authSVC authService.AuthService,
	validationSVC validationService.ValidationService,
) {
	authMiddleware := authMiddleware.New(authSVC)
	todoHandler := todoHandler.NewHandler(todoSVC, validationSVC)
	todoRoutes := v1.Group("/todo")

	todoRoutes.Use(authMiddleware.VerifyJWT())

	todoRoutes.POST("", todoHandler.CreateTodo)
	todoRoutes.GET("", todoHandler.GetTodos)
	todoRoutes.GET("/:todoId", todoHandler.GetTodo)
	todoRoutes.PATCH("/:todoId", todoHandler.MarkAsDone)
	todoRoutes.DELETE("/:todoId", todoHandler.DeleteTodo)
}
