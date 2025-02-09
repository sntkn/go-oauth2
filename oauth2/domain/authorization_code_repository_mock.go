// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package domain

import (
	"sync"
	"time"
)

// Ensure, that AuthorizationCodeRepositoryMock does implement AuthorizationCodeRepository.
// If this is not the case, regenerate this file with moq.
var _ AuthorizationCodeRepository = &AuthorizationCodeRepositoryMock{}

// AuthorizationCodeRepositoryMock is a mock implementation of AuthorizationCodeRepository.
//
//	func TestSomethingThatUsesAuthorizationCodeRepository(t *testing.T) {
//
//		// make and configure a mocked AuthorizationCodeRepository
//		mockedAuthorizationCodeRepository := &AuthorizationCodeRepositoryMock{
//			FindAuthorizationCodeFunc: func(s string) (AuthorizationCode, error) {
//				panic("mock out the FindAuthorizationCode method")
//			},
//			FindValidAuthorizationCodeFunc: func(s string, timeMoqParam time.Time) (AuthorizationCode, error) {
//				panic("mock out the FindValidAuthorizationCode method")
//			},
//			RevokeCodeFunc: func(code string) error {
//				panic("mock out the RevokeCode method")
//			},
//			StoreAuthorizationCodeFunc: func(storeAuthorizationCodeParams StoreAuthorizationCodeParams) (string, error) {
//				panic("mock out the StoreAuthorizationCode method")
//			},
//		}
//
//		// use mockedAuthorizationCodeRepository in code that requires AuthorizationCodeRepository
//		// and then make assertions.
//
//	}
type AuthorizationCodeRepositoryMock struct {
	// FindAuthorizationCodeFunc mocks the FindAuthorizationCode method.
	FindAuthorizationCodeFunc func(s string) (AuthorizationCode, error)

	// FindValidAuthorizationCodeFunc mocks the FindValidAuthorizationCode method.
	FindValidAuthorizationCodeFunc func(s string, timeMoqParam time.Time) (AuthorizationCode, error)

	// RevokeCodeFunc mocks the RevokeCode method.
	RevokeCodeFunc func(code string) error

	// StoreAuthorizationCodeFunc mocks the StoreAuthorizationCode method.
	StoreAuthorizationCodeFunc func(storeAuthorizationCodeParams StoreAuthorizationCodeParams) (string, error)

	// calls tracks calls to the methods.
	calls struct {
		// FindAuthorizationCode holds details about calls to the FindAuthorizationCode method.
		FindAuthorizationCode []struct {
			// S is the s argument value.
			S string
		}
		// FindValidAuthorizationCode holds details about calls to the FindValidAuthorizationCode method.
		FindValidAuthorizationCode []struct {
			// S is the s argument value.
			S string
			// TimeMoqParam is the timeMoqParam argument value.
			TimeMoqParam time.Time
		}
		// RevokeCode holds details about calls to the RevokeCode method.
		RevokeCode []struct {
			// Code is the code argument value.
			Code string
		}
		// StoreAuthorizationCode holds details about calls to the StoreAuthorizationCode method.
		StoreAuthorizationCode []struct {
			// StoreAuthorizationCodeParams is the storeAuthorizationCodeParams argument value.
			StoreAuthorizationCodeParams StoreAuthorizationCodeParams
		}
	}
	lockFindAuthorizationCode      sync.RWMutex
	lockFindValidAuthorizationCode sync.RWMutex
	lockRevokeCode                 sync.RWMutex
	lockStoreAuthorizationCode     sync.RWMutex
}

// FindAuthorizationCode calls FindAuthorizationCodeFunc.
func (mock *AuthorizationCodeRepositoryMock) FindAuthorizationCode(s string) (AuthorizationCode, error) {
	if mock.FindAuthorizationCodeFunc == nil {
		panic("AuthorizationCodeRepositoryMock.FindAuthorizationCodeFunc: method is nil but AuthorizationCodeRepository.FindAuthorizationCode was just called")
	}
	callInfo := struct {
		S string
	}{
		S: s,
	}
	mock.lockFindAuthorizationCode.Lock()
	mock.calls.FindAuthorizationCode = append(mock.calls.FindAuthorizationCode, callInfo)
	mock.lockFindAuthorizationCode.Unlock()
	return mock.FindAuthorizationCodeFunc(s)
}

// FindAuthorizationCodeCalls gets all the calls that were made to FindAuthorizationCode.
// Check the length with:
//
//	len(mockedAuthorizationCodeRepository.FindAuthorizationCodeCalls())
func (mock *AuthorizationCodeRepositoryMock) FindAuthorizationCodeCalls() []struct {
	S string
} {
	var calls []struct {
		S string
	}
	mock.lockFindAuthorizationCode.RLock()
	calls = mock.calls.FindAuthorizationCode
	mock.lockFindAuthorizationCode.RUnlock()
	return calls
}

// FindValidAuthorizationCode calls FindValidAuthorizationCodeFunc.
func (mock *AuthorizationCodeRepositoryMock) FindValidAuthorizationCode(s string, timeMoqParam time.Time) (AuthorizationCode, error) {
	if mock.FindValidAuthorizationCodeFunc == nil {
		panic("AuthorizationCodeRepositoryMock.FindValidAuthorizationCodeFunc: method is nil but AuthorizationCodeRepository.FindValidAuthorizationCode was just called")
	}
	callInfo := struct {
		S            string
		TimeMoqParam time.Time
	}{
		S:            s,
		TimeMoqParam: timeMoqParam,
	}
	mock.lockFindValidAuthorizationCode.Lock()
	mock.calls.FindValidAuthorizationCode = append(mock.calls.FindValidAuthorizationCode, callInfo)
	mock.lockFindValidAuthorizationCode.Unlock()
	return mock.FindValidAuthorizationCodeFunc(s, timeMoqParam)
}

// FindValidAuthorizationCodeCalls gets all the calls that were made to FindValidAuthorizationCode.
// Check the length with:
//
//	len(mockedAuthorizationCodeRepository.FindValidAuthorizationCodeCalls())
func (mock *AuthorizationCodeRepositoryMock) FindValidAuthorizationCodeCalls() []struct {
	S            string
	TimeMoqParam time.Time
} {
	var calls []struct {
		S            string
		TimeMoqParam time.Time
	}
	mock.lockFindValidAuthorizationCode.RLock()
	calls = mock.calls.FindValidAuthorizationCode
	mock.lockFindValidAuthorizationCode.RUnlock()
	return calls
}

// RevokeCode calls RevokeCodeFunc.
func (mock *AuthorizationCodeRepositoryMock) RevokeCode(code string) error {
	if mock.RevokeCodeFunc == nil {
		panic("AuthorizationCodeRepositoryMock.RevokeCodeFunc: method is nil but AuthorizationCodeRepository.RevokeCode was just called")
	}
	callInfo := struct {
		Code string
	}{
		Code: code,
	}
	mock.lockRevokeCode.Lock()
	mock.calls.RevokeCode = append(mock.calls.RevokeCode, callInfo)
	mock.lockRevokeCode.Unlock()
	return mock.RevokeCodeFunc(code)
}

// RevokeCodeCalls gets all the calls that were made to RevokeCode.
// Check the length with:
//
//	len(mockedAuthorizationCodeRepository.RevokeCodeCalls())
func (mock *AuthorizationCodeRepositoryMock) RevokeCodeCalls() []struct {
	Code string
} {
	var calls []struct {
		Code string
	}
	mock.lockRevokeCode.RLock()
	calls = mock.calls.RevokeCode
	mock.lockRevokeCode.RUnlock()
	return calls
}

// StoreAuthorizationCode calls StoreAuthorizationCodeFunc.
func (mock *AuthorizationCodeRepositoryMock) StoreAuthorizationCode(storeAuthorizationCodeParams StoreAuthorizationCodeParams) (string, error) {
	if mock.StoreAuthorizationCodeFunc == nil {
		panic("AuthorizationCodeRepositoryMock.StoreAuthorizationCodeFunc: method is nil but AuthorizationCodeRepository.StoreAuthorizationCode was just called")
	}
	callInfo := struct {
		StoreAuthorizationCodeParams StoreAuthorizationCodeParams
	}{
		StoreAuthorizationCodeParams: storeAuthorizationCodeParams,
	}
	mock.lockStoreAuthorizationCode.Lock()
	mock.calls.StoreAuthorizationCode = append(mock.calls.StoreAuthorizationCode, callInfo)
	mock.lockStoreAuthorizationCode.Unlock()
	return mock.StoreAuthorizationCodeFunc(storeAuthorizationCodeParams)
}

// StoreAuthorizationCodeCalls gets all the calls that were made to StoreAuthorizationCode.
// Check the length with:
//
//	len(mockedAuthorizationCodeRepository.StoreAuthorizationCodeCalls())
func (mock *AuthorizationCodeRepositoryMock) StoreAuthorizationCodeCalls() []struct {
	StoreAuthorizationCodeParams StoreAuthorizationCodeParams
} {
	var calls []struct {
		StoreAuthorizationCodeParams StoreAuthorizationCodeParams
	}
	mock.lockStoreAuthorizationCode.RLock()
	calls = mock.calls.StoreAuthorizationCode
	mock.lockStoreAuthorizationCode.RUnlock()
	return calls
}
