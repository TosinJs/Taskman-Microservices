package responseEntity

import "tosinjs/notification-service/internal/entity/errorEntity"

type ResponseEntity struct {
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
	Error   any    `json:"error,omitempty"`
	Path    string `json:"path,omitempty"`
}

func BuildResponseObject(path string, body any) ResponseEntity {
	return ResponseEntity{
		Message: "success",
		Data:    body,
		Path:    path,
	}
}

func BuildErrorResponseObject(errorMessage, path string) ResponseEntity {
	return ResponseEntity{
		Message: "failed",
		Error:   errorMessage,
		Path:    path,
	}
}

func BuildServiceErrorResponseObject(err *errorEntity.LayerError, path string) ResponseEntity {
	return ResponseEntity{
		Message: "failed",
		Error:   err.Error(),
		Path:    path,
	}
}
