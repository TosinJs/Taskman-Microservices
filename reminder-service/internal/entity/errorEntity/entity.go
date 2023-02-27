package errorEntity

import (
	"fmt"
	"net/http"
)

type LayerError struct {
	Layer      string
	StatusCode int
	message    string
	err        error
}

func (se LayerError) Error() string {
	return fmt.Sprintf("Error: %s", se.message)
}

func New(statusCode int, layer, message string, err error) *LayerError {
	return &LayerError{
		Layer:      layer,
		StatusCode: statusCode,
		message:    message,
		err:        err,
	}
}

func ConflictError(layer, message string, err error) *LayerError {
	return New(http.StatusConflict, layer, message, err)
}

func InternalServerError(layer string, err error) *LayerError {
	return New(http.StatusInternalServerError, layer, "Internal Server Error", err)
}

func NotFoundError(layer, message string, err error) *LayerError {
	return New(http.StatusNotFound, layer, message, err)
}

func BadRequestError(layer, message string, err error) *LayerError {
	return New(http.StatusBadRequest, layer, message, err)
}
