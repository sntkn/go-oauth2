// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package valkey

import (
	"context"
	"sync"
)

// Ensure, that ClientIFMock does implement ClientIF.
// If this is not the case, regenerate this file with moq.
var _ ClientIF = &ClientIFMock{}

// ClientIFMock is a mock implementation of ClientIF.
//
//	func TestSomethingThatUsesClientIF(t *testing.T) {
//
//		// make and configure a mocked ClientIF
//		mockedClientIF := &ClientIFMock{
//			DelFunc: func(ctx context.Context, key string) error {
//				panic("mock out the Del method")
//			},
//			GetFunc: func(ctx context.Context, key string) (string, error) {
//				panic("mock out the Get method")
//			},
//			SetFunc: func(ctx context.Context, key string, value string, expiration int64) error {
//				panic("mock out the Set method")
//			},
//		}
//
//		// use mockedClientIF in code that requires ClientIF
//		// and then make assertions.
//
//	}
type ClientIFMock struct {
	// DelFunc mocks the Del method.
	DelFunc func(ctx context.Context, key string) error

	// GetFunc mocks the Get method.
	GetFunc func(ctx context.Context, key string) (string, error)

	// SetFunc mocks the Set method.
	SetFunc func(ctx context.Context, key string, value string, expiration int64) error

	// calls tracks calls to the methods.
	calls struct {
		// Del holds details about calls to the Del method.
		Del []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Key is the key argument value.
			Key string
		}
		// Get holds details about calls to the Get method.
		Get []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Key is the key argument value.
			Key string
		}
		// Set holds details about calls to the Set method.
		Set []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Key is the key argument value.
			Key string
			// Value is the value argument value.
			Value string
			// Expiration is the expiration argument value.
			Expiration int64
		}
	}
	lockDel sync.RWMutex
	lockGet sync.RWMutex
	lockSet sync.RWMutex
}

// Del calls DelFunc.
func (mock *ClientIFMock) Del(ctx context.Context, key string) error {
	if mock.DelFunc == nil {
		panic("ClientIFMock.DelFunc: method is nil but ClientIF.Del was just called")
	}
	callInfo := struct {
		Ctx context.Context
		Key string
	}{
		Ctx: ctx,
		Key: key,
	}
	mock.lockDel.Lock()
	mock.calls.Del = append(mock.calls.Del, callInfo)
	mock.lockDel.Unlock()
	return mock.DelFunc(ctx, key)
}

// DelCalls gets all the calls that were made to Del.
// Check the length with:
//
//	len(mockedClientIF.DelCalls())
func (mock *ClientIFMock) DelCalls() []struct {
	Ctx context.Context
	Key string
} {
	var calls []struct {
		Ctx context.Context
		Key string
	}
	mock.lockDel.RLock()
	calls = mock.calls.Del
	mock.lockDel.RUnlock()
	return calls
}

// Get calls GetFunc.
func (mock *ClientIFMock) Get(ctx context.Context, key string) (string, error) {
	if mock.GetFunc == nil {
		panic("ClientIFMock.GetFunc: method is nil but ClientIF.Get was just called")
	}
	callInfo := struct {
		Ctx context.Context
		Key string
	}{
		Ctx: ctx,
		Key: key,
	}
	mock.lockGet.Lock()
	mock.calls.Get = append(mock.calls.Get, callInfo)
	mock.lockGet.Unlock()
	return mock.GetFunc(ctx, key)
}

// GetCalls gets all the calls that were made to Get.
// Check the length with:
//
//	len(mockedClientIF.GetCalls())
func (mock *ClientIFMock) GetCalls() []struct {
	Ctx context.Context
	Key string
} {
	var calls []struct {
		Ctx context.Context
		Key string
	}
	mock.lockGet.RLock()
	calls = mock.calls.Get
	mock.lockGet.RUnlock()
	return calls
}

// Set calls SetFunc.
func (mock *ClientIFMock) Set(ctx context.Context, key string, value string, expiration int64) error {
	if mock.SetFunc == nil {
		panic("ClientIFMock.SetFunc: method is nil but ClientIF.Set was just called")
	}
	callInfo := struct {
		Ctx        context.Context
		Key        string
		Value      string
		Expiration int64
	}{
		Ctx:        ctx,
		Key:        key,
		Value:      value,
		Expiration: expiration,
	}
	mock.lockSet.Lock()
	mock.calls.Set = append(mock.calls.Set, callInfo)
	mock.lockSet.Unlock()
	return mock.SetFunc(ctx, key, value, expiration)
}

// SetCalls gets all the calls that were made to Set.
// Check the length with:
//
//	len(mockedClientIF.SetCalls())
func (mock *ClientIFMock) SetCalls() []struct {
	Ctx        context.Context
	Key        string
	Value      string
	Expiration int64
} {
	var calls []struct {
		Ctx        context.Context
		Key        string
		Value      string
		Expiration int64
	}
	mock.lockSet.RLock()
	calls = mock.calls.Set
	mock.lockSet.RUnlock()
	return calls
}
