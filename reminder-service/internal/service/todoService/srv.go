package todoService

import (
	"fmt"
	"time"
	"tosinjs/reminder-service/internal/entity/errorEntity"
	"tosinjs/reminder-service/internal/entity/todoEntity"
	"tosinjs/reminder-service/internal/repository/todoRepo"
	"tosinjs/reminder-service/internal/service/reminderService"
)

type todoService struct {
	todoRepo    todoRepo.TodoRepository
	reminderSVC reminderService.ReminderService
}

type TodoService interface {
	CreateTodo(todo todoEntity.CreateTodoReq) *errorEntity.LayerError
	GetTodos(userId string) ([]todoEntity.Todo, *errorEntity.LayerError)
	GetTodo(userId, todoId string) (todoEntity.Todo, *errorEntity.LayerError)
	DeleteTodo(userId, todoId string) *errorEntity.LayerError
	MarkAsDone(userId, todoId string) *errorEntity.LayerError
}

// Using RFC 3339 Time Format
const timeLayout = "2006-01-02T15:04:05Z"

func New(
	todoRepo todoRepo.TodoRepository,
	reminderSVC reminderService.ReminderService,
) TodoService {
	return todoService{
		todoRepo:    todoRepo,
		reminderSVC: reminderSVC,
	}
}

func (t todoService) CreateTodo(req todoEntity.CreateTodoReq) *errorEntity.LayerError {
	now := time.Now().UTC()
	due, err := time.Parse(timeLayout, req.Due)
	if err != nil {
		return errorEntity.BadRequestError(
			"service", "Inavlid Time Input: Due", err,
		)
	}
	if int(due.Sub(now).Minutes()) < 2 {
		return errorEntity.BadRequestError(
			"service", "due has to be due at least 1 minute from now", fmt.Errorf("inavlid time entry"),
		)
	}

	var remindMe time.Time
	if req.RemindMe != "" {
		remindMe, err = time.Parse(timeLayout, req.RemindMe)
		if err != nil {
			return errorEntity.BadRequestError(
				"service", "Inavlid Time Input: RemindMe", err,
			)
		}
		if int(remindMe.Sub(now).Minutes()) < 2 {
			return errorEntity.BadRequestError(
				"service",
				"remindMe has to be due at least 1 minute from now",
				fmt.Errorf("inavlid time entry"),
			)
		}
	}

	todo := todoEntity.Todo{
		UserId:    req.UserId,
		CreatedAt: time.Now(),
		Due:       due,
		Todo:      req.Todo,
		RemindMe:  remindMe,
		Recurring: req.Recurring,
		Done:      false,
	}

	todoId, svcErr := t.todoRepo.CreateTodo(todo)
	if err != nil {
		return svcErr
	}

	//Create a reminder for when the task is due
	t.reminderSVC.CreateReminder(
		req.UserId,
		todoId,
		"This task is expiring now",
		due,
	)

	//Create a reminder based on the remindMe
	if !remindMe.IsZero() {
		fmt.Println(remindMe)
		t.reminderSVC.CreateReminder(
			req.UserId,
			todoId,
			"This is your preset reminder to complete this task",
			remindMe,
		)
	}
	return nil
}

func (t todoService) GetTodos(userId string) ([]todoEntity.Todo, *errorEntity.LayerError) {
	return t.todoRepo.GetTodos(userId)
}

func (t todoService) GetTodo(userId, todoId string) (todoEntity.Todo, *errorEntity.LayerError) {
	return t.todoRepo.GetTodo(userId, todoId)
}

func (t todoService) DeleteTodo(userId, todoId string) *errorEntity.LayerError {
	//If the todo is deleted the reminders should go to
	t.reminderSVC.DeleteReminder(todoId)
	return t.todoRepo.DeleteTodo(userId, todoId)
}
func (t todoService) MarkAsDone(userId, todoId string) *errorEntity.LayerError {
	//If the todo is done there is nothing to remind you about
	t.reminderSVC.DeleteReminder(todoId)
	return t.todoRepo.MarkAsDone(userId, todoId)
}
