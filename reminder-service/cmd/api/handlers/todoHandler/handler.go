package todoHandler

import (
	"net/http"
	"tosinjs/reminder-service/internal/entity/responseEntity"
	"tosinjs/reminder-service/internal/entity/todoEntity"
	"tosinjs/reminder-service/internal/service/todoService"
	"tosinjs/reminder-service/internal/service/validationService"

	"github.com/gin-gonic/gin"
)

type todoHandler struct {
	todoSVC       todoService.TodoService
	validationSVC validationService.ValidationService
}

func NewHandler(
	todoSVC todoService.TodoService,
	validationSVC validationService.ValidationService,
) todoHandler {
	return todoHandler{
		todoSVC:       todoSVC,
		validationSVC: validationSVC,
	}
}

func (t todoHandler) CreateTodo(c *gin.Context) {
	userId := c.GetString("userId")
	var req todoEntity.CreateTodoReq
	req.UserId = userId

	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, responseEntity.BuildErrorResponseObject(
			err.Error(), c.FullPath(),
		))
		return
	}
	if svcErr := t.validationSVC.Validate(req); svcErr != nil {
		c.AbortWithStatusJSON(svcErr.StatusCode, responseEntity.BuildServiceErrorResponseObject(
			svcErr, c.FullPath(),
		))
		return
	}
	if svcErr := t.todoSVC.CreateTodo(req); svcErr != nil {
		c.AbortWithStatusJSON(svcErr.StatusCode, responseEntity.BuildServiceErrorResponseObject(
			svcErr, c.FullPath(),
		))
		return
	}
	c.JSON(http.StatusAccepted, "success")
}

func (t todoHandler) GetTodos(c *gin.Context) {
	userId := c.GetString("userId")

	todos, svcErr := t.todoSVC.GetTodos(userId)
	if svcErr != nil {
		c.AbortWithStatusJSON(svcErr.StatusCode, responseEntity.BuildServiceErrorResponseObject(
			svcErr, c.FullPath(),
		))
		return
	}

	c.JSON(http.StatusAccepted, todos)
}
func (t todoHandler) GetTodo(c *gin.Context) {
	userId := c.GetString("userId")
	todoId := c.Params.ByName("todoId")

	todo, svcErr := t.todoSVC.GetTodo(userId, todoId)
	if svcErr != nil {
		c.AbortWithStatusJSON(svcErr.StatusCode, responseEntity.BuildServiceErrorResponseObject(
			svcErr, c.FullPath(),
		))
		return
	}
	c.JSON(http.StatusAccepted, todo)
}
func (t todoHandler) DeleteTodo(c *gin.Context) {
	userId := c.GetString("userId")
	todoId := c.Params.ByName("todoId")

	if svcErr := t.todoSVC.DeleteTodo(userId, todoId); svcErr != nil {
		c.AbortWithStatusJSON(svcErr.StatusCode, responseEntity.BuildServiceErrorResponseObject(
			svcErr, c.FullPath(),
		))
		return
	}
	c.JSON(http.StatusAccepted, "deleted")
}
func (t todoHandler) MarkAsDone(c *gin.Context) {
	userId := c.GetString("userId")
	todoId := c.Params.ByName("todoId")
	if svcErr := t.todoSVC.MarkAsDone(userId, todoId); svcErr != nil {
		c.AbortWithStatusJSON(svcErr.StatusCode, responseEntity.BuildServiceErrorResponseObject(
			svcErr, c.FullPath(),
		))
		return
	}
	c.JSON(http.StatusAccepted, "updated")
}
