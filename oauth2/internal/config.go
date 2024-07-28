package internal

import (
	"fmt"
	"reflect"

	"github.com/gin-gonic/gin"
)

func GetFromContextIF[T any](c *gin.Context, name string) (T, error) {
	var zero T

	// コンテキストから値を取得
	s, exists := c.Get(name)
	if !exists {
		return zero, fmt.Errorf("context value '%s' not found", name)
	}

	// 期待する型を取得
	expectedType := reflect.TypeOf((*T)(nil)).Elem()
	if expectedType.Kind() != reflect.Interface {
		return zero, fmt.Errorf("expected type '%s' is not an interface", expectedType.Name())
	}

	// 実際の値の型を取得
	valueType := reflect.TypeOf(s)
	if !valueType.Implements(expectedType) {
		return zero, fmt.Errorf("context value '%s' does not implement expected type '%s', got '%s'", name, expectedType.Name(), valueType.Name())
	}

	// 型アサーションのチェック
	castedValue, ok := s.(T)
	if !ok {
		return zero, fmt.Errorf("failed to cast context value '%s' to type '%s'", name, expectedType.Name())
	}

	return castedValue, nil
}

func GetFromContext[T any](c *gin.Context, name string) (*T, error) {
	// コンテキストから値を取得
	s, exists := c.Get(name)
	if !exists {
		return nil, fmt.Errorf("context value '%s' not found", name)
	}

	// 期待する型のインスタンスを取得
	var zero T
	expectedType := reflect.TypeOf(&zero).Elem()
	actualType := reflect.TypeOf(s)

	// 実際の型が期待する型でない場合はエラー
	if actualType.Kind() == reflect.Ptr && actualType.Elem() == expectedType {
		// 型アサーションが成功するか確認
		castedValue, ok := s.(*T)
		if !ok {
			return nil, fmt.Errorf("context value '%s' cannot be cast to type '%s'", name, expectedType.Name())
		}
		return castedValue, nil
	}

	// 期待する型がインターフェースの場合はエラー
	if expectedType.Kind() == reflect.Interface {
		return nil, fmt.Errorf("expected type '%s' is an interface, but value in context is not", expectedType.Name())
	}

	return nil, fmt.Errorf("context value '%s' is not of expected type '%s', got '%s'", name, expectedType.Name(), actualType.Name())
}
