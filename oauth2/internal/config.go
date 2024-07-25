package internal

import (
	"fmt"
	"reflect"

	"github.com/gin-gonic/gin"
)

func GetFromContext[T any](c *gin.Context, name string) (*T, error) {
	s, exists := c.Get(name)
	if !exists {
		return nil, fmt.Errorf("%s not found", name)
	}

	sessionValue, ok := s.(*T)
	if !ok {
		return nil, fmt.Errorf("%s value is not of expected type %s, got %s", name, reflect.TypeOf((*T)(nil)).Elem().Name(), reflect.TypeOf(s).Name())

	}

	return sessionValue, nil
}
