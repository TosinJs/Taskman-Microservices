package validationService

import (
	"tosinjs/reminder-service/internal/entity/errorEntity"

	"github.com/go-playground/validator/v10"
)

type validationService struct{}

type ValidationService interface {
	Validate(value any) *errorEntity.LayerError
}

func New() ValidationService {
	return validationService{}
}

func (v validationService) Validate(value any) *errorEntity.LayerError {
	validate := validator.New()
	if svcErr := validate.Struct(value); svcErr != nil {
		return errorEntity.BadRequestError("service", svcErr.Error(), svcErr)
	}
	return nil
}
