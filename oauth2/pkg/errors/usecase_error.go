package errors

import "fmt"

type UsecaseError struct {
	Code        int
	Message     string
	RedirectURI string
}

func (e *UsecaseError) Error() string {
	return fmt.Sprintf("Code: %d, Message: %s", e.Code, e.Message)
}

func NewUsecaseError(code int, message string) *UsecaseError {
	return &UsecaseError{
		Code:    code,
		Message: message,
	}
}

func NewUsecaseErrorWithRedirectURI(code int, message, redirectURI string) *UsecaseError {
	return &UsecaseError{
		Code:        code,
		Message:     message,
		RedirectURI: redirectURI,
	}
}
