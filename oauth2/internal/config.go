package internal

import (
	"fmt"
	"reflect"

	"github.com/gin-gonic/gin"
)

func GetFromContextIF[T any](c *gin.Context, name string) (T, error) {
	var zero T
	s, exists := c.Get(name)
	if !exists {
		return zero, fmt.Errorf("%s not found", name)
	}

	// 期待する型を取得
	expectedType := reflect.TypeOf((*T)(nil)).Elem()
	valueType := reflect.TypeOf(s)

	// 実際の型がインターフェース型であるかの確認
	if expectedType.Kind() != reflect.Interface {
		return zero, fmt.Errorf("%s is not an interface", expectedType.Name())
	}

	// 期待する型がインターフェースであるかを確認
	if !expectedType.Implements(reflect.TypeOf((*interface{})(nil)).Elem()) {
		return zero, fmt.Errorf("%s is not an interface", expectedType.Name())
	}

	// 実際の型が期待するインターフェースを実装しているか確認
	if !valueType.Implements(expectedType) {
		return zero, fmt.Errorf("%s value is not of expected type %s, got %s", name, expectedType.Name(), valueType.Name())
	}

	castedValue, ok := s.(T)
	if !ok {
		return zero, fmt.Errorf("failed to cast %s to %s", valueType.Name(), expectedType.Name())
	}

	return castedValue, nil
}

func GetFromContext[T any](c *gin.Context, name string) (*T, error) {
	s, exists := c.Get(name)
	if !exists {
		return nil, fmt.Errorf("%s not found", name)
	}

	// ポインタのインターフェース型へのアサーション
	var t T
	interfaceType := reflect.TypeOf(&t).Elem()
	valueType := reflect.TypeOf(s)

	if interfaceType.Kind() == reflect.Interface {
		return nil, fmt.Errorf("%s value is interface", valueType.Elem().Name())
	}
	value, ok := s.(*T)
	if !ok {
		return nil, fmt.Errorf("%s value is not of expected type %s, got %s", name, reflect.TypeOf((*T)(nil)).Name(), reflect.TypeOf(s).Elem().Name())
	}
	return value, nil
}
