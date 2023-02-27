package todoRepo

import (
	"tosinjs/reminder-service/internal/entity/errorEntity"
	"tosinjs/reminder-service/internal/entity/todoEntity"
)

type TodoRepository interface {
	CreateTodo(todo todoEntity.Todo) (string, *errorEntity.LayerError)
	GetTodos(userId string) ([]todoEntity.Todo, *errorEntity.LayerError)
	GetTodo(userId, todoId string) (todoEntity.Todo, *errorEntity.LayerError)
	DeleteTodo(userId, todoId string) *errorEntity.LayerError
	MarkAsDone(userId, todoId string) *errorEntity.LayerError
}
