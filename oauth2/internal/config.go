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

	// ポインタのインターフェース型へのアサーション
	var t T
	interfaceType := reflect.TypeOf(&t).Elem()
	valueType := reflect.TypeOf(s)

	if interfaceType.Kind() != reflect.Interface {
		value, ok := s.(*T)
		if !ok {
			return nil, fmt.Errorf("%s value is not of expected type %s, got %s", name, reflect.TypeOf((*T)(nil)).Name(), reflect.TypeOf(s).Elem().Name())
		}
		return value, nil
	}

	if !valueType.Implements(interfaceType) {
		return nil, fmt.Errorf("%s value is not of expected type %s, got %s", name, interfaceType.Name(), valueType.Name())
	}

	castedValue := s.(T)
	return &castedValue, nil
}
