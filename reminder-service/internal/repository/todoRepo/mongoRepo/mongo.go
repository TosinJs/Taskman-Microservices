package mongoRepo

import (
	"context"
	"tosinjs/reminder-service/internal/entity/errorEntity"
	"tosinjs/reminder-service/internal/entity/todoEntity"
	"tosinjs/reminder-service/internal/repository/todoRepo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoRepo struct {
	mColl *mongo.Collection
	ctx   context.Context
}

func New(mColl *mongo.Collection, ctx context.Context) todoRepo.TodoRepository {
	return mongoRepo{
		mColl: mColl,
		ctx:   ctx,
	}
}

func (m mongoRepo) CreateTodo(todo todoEntity.Todo) (string, *errorEntity.LayerError) {
	todo.Id = primitive.NewObjectID()
	_, err := m.mColl.InsertOne(m.ctx, todo)
	return todo.Id.String(), errorEntity.InternalServerError("repo", err)
}

func (m mongoRepo) GetTodos(userId string) ([]todoEntity.Todo, *errorEntity.LayerError) {
	filter := bson.D{{Key: "userId", Value: userId}}
	todoCursor, err := m.mColl.Find(m.ctx, filter)
	if err != nil {
		return nil, errorEntity.InternalServerError("repo", err)
	}

	todos := make([]todoEntity.Todo, 0)
	for todoCursor.Next(m.ctx) {
		var todo todoEntity.Todo
		err = todoCursor.Decode(&todo)
		todos = append(todos, todo)
	}

	if err != nil {
		return nil, errorEntity.InternalServerError("repo", err)
	}
	return todos, nil
}

func (m mongoRepo) GetTodo(userId, todoId string) (todoEntity.Todo, *errorEntity.LayerError) {
	todoIdHex, err := primitive.ObjectIDFromHex(todoId)
	if err != nil {
		return todoEntity.Todo{}, errorEntity.BadRequestError("repo", "Invalid TodoId", err)
	}
	filter := bson.M{"_id": todoIdHex, "userId": userId}
	todoRes := m.mColl.FindOne(m.ctx, filter)
	if todoRes.Err() == nil {
		if todoRes.Err() == mongo.ErrNoDocuments {
			return todoEntity.Todo{}, errorEntity.NotFoundError(
				"repo", "Todo Not Found", todoRes.Err(),
			)
		}
		return todoEntity.Todo{}, errorEntity.InternalServerError("repo", todoRes.Err())
	}

	var todo todoEntity.Todo
	err = todoRes.Decode(&todo)
	if err != nil {
		return todoEntity.Todo{}, errorEntity.InternalServerError("repo", todoRes.Err())
	}

	return todo, nil
}

func (m mongoRepo) DeleteTodo(userId, todoId string) *errorEntity.LayerError {
	todoIdHex, err := primitive.ObjectIDFromHex(todoId)
	if err != nil {
		return errorEntity.BadRequestError("repo", "Invalid TodoId", err)
	}
	filter := bson.M{"_id": todoIdHex, "userId": userId}
	_, err = m.mColl.DeleteOne(m.ctx, filter)
	return errorEntity.InternalServerError("repo", err)
}

func (m mongoRepo) MarkAsDone(userId, todoId string) *errorEntity.LayerError {
	todoIdHex, err := primitive.ObjectIDFromHex(todoId)
	if err != nil {
		return errorEntity.BadRequestError("repo", "Invalid TodoId", err)
	}
	filter := bson.M{"_id": todoIdHex, "userId": userId}
	_, err = m.mColl.UpdateOne(m.ctx, filter, bson.M{"$set": bson.M{"done": true}})
	if err != nil {
		return errorEntity.InternalServerError("repo", err)
	}
	return nil
}
