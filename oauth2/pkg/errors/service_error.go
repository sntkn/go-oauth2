package errors

import "fmt"

const (
	// ErrCodeInternalServer is a code for internal server error
	ErrCodeInternalServer = 500
	ErrCodeForbidden      = 403
	ErrCodeNotFound       = 404
)

type ServiceError struct {
	Code        int
	Message     string
	RedirectURI string
}

func (e *ServiceError) Error() string {
	return fmt.Sprintf("Code: %d, Message: %s", e.Code, e.Message)
}

func NewServiceErrorError(code int, message string) *ServiceError {
	return &ServiceError{
		Code:    code,
		Message: message,
	}
}
