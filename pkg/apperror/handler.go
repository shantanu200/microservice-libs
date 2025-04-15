package apperror

import "github.com/gofiber/fiber/v2"

type AppError struct {
	Code    int
	Message string
}

func (r *AppError) Error() string {
	return r.Message
}

func NewAppError(code int, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

func InternalServerError(message string) *AppError {
	return NewAppError(fiber.StatusInternalServerError, message)
}

func UnprocessableEntityError(message string) *AppError {
	return NewAppError(fiber.StatusUnprocessableEntity, message)
}
