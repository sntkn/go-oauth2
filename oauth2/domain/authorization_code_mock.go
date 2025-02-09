// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package domain

import (
	"github.com/google/uuid"
	"sync"
	"time"
)

// Ensure, that AuthorizationCodeMock does implement AuthorizationCode.
// If this is not the case, regenerate this file with moq.
var _ AuthorizationCode = &AuthorizationCodeMock{}

// AuthorizationCodeMock is a mock implementation of AuthorizationCode.
//
//	func TestSomethingThatUsesAuthorizationCode(t *testing.T) {
//
//		// make and configure a mocked AuthorizationCode
//		mockedAuthorizationCode := &AuthorizationCodeMock{
//			GenerateRedirectURIWithCodeFunc: func() string {
//				panic("mock out the GenerateRedirectURIWithCode method")
//			},
//			GetClientIDFunc: func() uuid.UUID {
//				panic("mock out the GetClientID method")
//			},
//			GetCodeFunc: func() string {
//				panic("mock out the GetCode method")
//			},
//			GetExpiresAtFunc: func() time.Time {
//				panic("mock out the GetExpiresAt method")
//			},
//			GetRedirectURIFunc: func() string {
//				panic("mock out the GetRedirectURI method")
//			},
//			GetScopeFunc: func() string {
//				panic("mock out the GetScope method")
//			},
//			GetUserIDFunc: func() uuid.UUID {
//				panic("mock out the GetUserID method")
//			},
//			IsExpiredFunc: func(t time.Time) bool {
//				panic("mock out the IsExpired method")
//			},
//			IsNotFoundFunc: func() bool {
//				panic("mock out the IsNotFound method")
//			},
//		}
//
//		// use mockedAuthorizationCode in code that requires AuthorizationCode
//		// and then make assertions.
//
//	}
type AuthorizationCodeMock struct {
	// GenerateRedirectURIWithCodeFunc mocks the GenerateRedirectURIWithCode method.
	GenerateRedirectURIWithCodeFunc func() string

	// GetClientIDFunc mocks the GetClientID method.
	GetClientIDFunc func() uuid.UUID

	// GetCodeFunc mocks the GetCode method.
	GetCodeFunc func() string

	// GetExpiresAtFunc mocks the GetExpiresAt method.
	GetExpiresAtFunc func() time.Time

	// GetRedirectURIFunc mocks the GetRedirectURI method.
	GetRedirectURIFunc func() string

	// GetScopeFunc mocks the GetScope method.
	GetScopeFunc func() string

	// GetUserIDFunc mocks the GetUserID method.
	GetUserIDFunc func() uuid.UUID

	// IsExpiredFunc mocks the IsExpired method.
	IsExpiredFunc func(t time.Time) bool

	// IsNotFoundFunc mocks the IsNotFound method.
	IsNotFoundFunc func() bool

	// calls tracks calls to the methods.
	calls struct {
		// GenerateRedirectURIWithCode holds details about calls to the GenerateRedirectURIWithCode method.
		GenerateRedirectURIWithCode []struct {
		}
		// GetClientID holds details about calls to the GetClientID method.
		GetClientID []struct {
		}
		// GetCode holds details about calls to the GetCode method.
		GetCode []struct {
		}
		// GetExpiresAt holds details about calls to the GetExpiresAt method.
		GetExpiresAt []struct {
		}
		// GetRedirectURI holds details about calls to the GetRedirectURI method.
		GetRedirectURI []struct {
		}
		// GetScope holds details about calls to the GetScope method.
		GetScope []struct {
		}
		// GetUserID holds details about calls to the GetUserID method.
		GetUserID []struct {
		}
		// IsExpired holds details about calls to the IsExpired method.
		IsExpired []struct {
			// T is the t argument value.
			T time.Time
		}
		// IsNotFound holds details about calls to the IsNotFound method.
		IsNotFound []struct {
		}
	}
	lockGenerateRedirectURIWithCode sync.RWMutex
	lockGetClientID                 sync.RWMutex
	lockGetCode                     sync.RWMutex
	lockGetExpiresAt                sync.RWMutex
	lockGetRedirectURI              sync.RWMutex
	lockGetScope                    sync.RWMutex
	lockGetUserID                   sync.RWMutex
	lockIsExpired                   sync.RWMutex
	lockIsNotFound                  sync.RWMutex
}

// GenerateRedirectURIWithCode calls GenerateRedirectURIWithCodeFunc.
func (mock *AuthorizationCodeMock) GenerateRedirectURIWithCode() string {
	if mock.GenerateRedirectURIWithCodeFunc == nil {
		panic("AuthorizationCodeMock.GenerateRedirectURIWithCodeFunc: method is nil but AuthorizationCode.GenerateRedirectURIWithCode was just called")
	}
	callInfo := struct {
	}{}
	mock.lockGenerateRedirectURIWithCode.Lock()
	mock.calls.GenerateRedirectURIWithCode = append(mock.calls.GenerateRedirectURIWithCode, callInfo)
	mock.lockGenerateRedirectURIWithCode.Unlock()
	return mock.GenerateRedirectURIWithCodeFunc()
}

// GenerateRedirectURIWithCodeCalls gets all the calls that were made to GenerateRedirectURIWithCode.
// Check the length with:
//
//	len(mockedAuthorizationCode.GenerateRedirectURIWithCodeCalls())
func (mock *AuthorizationCodeMock) GenerateRedirectURIWithCodeCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockGenerateRedirectURIWithCode.RLock()
	calls = mock.calls.GenerateRedirectURIWithCode
	mock.lockGenerateRedirectURIWithCode.RUnlock()
	return calls
}

// GetClientID calls GetClientIDFunc.
func (mock *AuthorizationCodeMock) GetClientID() uuid.UUID {
	if mock.GetClientIDFunc == nil {
		panic("AuthorizationCodeMock.GetClientIDFunc: method is nil but AuthorizationCode.GetClientID was just called")
	}
	callInfo := struct {
	}{}
	mock.lockGetClientID.Lock()
	mock.calls.GetClientID = append(mock.calls.GetClientID, callInfo)
	mock.lockGetClientID.Unlock()
	return mock.GetClientIDFunc()
}

// GetClientIDCalls gets all the calls that were made to GetClientID.
// Check the length with:
//
//	len(mockedAuthorizationCode.GetClientIDCalls())
func (mock *AuthorizationCodeMock) GetClientIDCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockGetClientID.RLock()
	calls = mock.calls.GetClientID
	mock.lockGetClientID.RUnlock()
	return calls
}

// GetCode calls GetCodeFunc.
func (mock *AuthorizationCodeMock) GetCode() string {
	if mock.GetCodeFunc == nil {
		panic("AuthorizationCodeMock.GetCodeFunc: method is nil but AuthorizationCode.GetCode was just called")
	}
	callInfo := struct {
	}{}
	mock.lockGetCode.Lock()
	mock.calls.GetCode = append(mock.calls.GetCode, callInfo)
	mock.lockGetCode.Unlock()
	return mock.GetCodeFunc()
}

// GetCodeCalls gets all the calls that were made to GetCode.
// Check the length with:
//
//	len(mockedAuthorizationCode.GetCodeCalls())
func (mock *AuthorizationCodeMock) GetCodeCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockGetCode.RLock()
	calls = mock.calls.GetCode
	mock.lockGetCode.RUnlock()
	return calls
}

// GetExpiresAt calls GetExpiresAtFunc.
func (mock *AuthorizationCodeMock) GetExpiresAt() time.Time {
	if mock.GetExpiresAtFunc == nil {
		panic("AuthorizationCodeMock.GetExpiresAtFunc: method is nil but AuthorizationCode.GetExpiresAt was just called")
	}
	callInfo := struct {
	}{}
	mock.lockGetExpiresAt.Lock()
	mock.calls.GetExpiresAt = append(mock.calls.GetExpiresAt, callInfo)
	mock.lockGetExpiresAt.Unlock()
	return mock.GetExpiresAtFunc()
}

// GetExpiresAtCalls gets all the calls that were made to GetExpiresAt.
// Check the length with:
//
//	len(mockedAuthorizationCode.GetExpiresAtCalls())
func (mock *AuthorizationCodeMock) GetExpiresAtCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockGetExpiresAt.RLock()
	calls = mock.calls.GetExpiresAt
	mock.lockGetExpiresAt.RUnlock()
	return calls
}

// GetRedirectURI calls GetRedirectURIFunc.
func (mock *AuthorizationCodeMock) GetRedirectURI() string {
	if mock.GetRedirectURIFunc == nil {
		panic("AuthorizationCodeMock.GetRedirectURIFunc: method is nil but AuthorizationCode.GetRedirectURI was just called")
	}
	callInfo := struct {
	}{}
	mock.lockGetRedirectURI.Lock()
	mock.calls.GetRedirectURI = append(mock.calls.GetRedirectURI, callInfo)
	mock.lockGetRedirectURI.Unlock()
	return mock.GetRedirectURIFunc()
}

// GetRedirectURICalls gets all the calls that were made to GetRedirectURI.
// Check the length with:
//
//	len(mockedAuthorizationCode.GetRedirectURICalls())
func (mock *AuthorizationCodeMock) GetRedirectURICalls() []struct {
} {
	var calls []struct {
	}
	mock.lockGetRedirectURI.RLock()
	calls = mock.calls.GetRedirectURI
	mock.lockGetRedirectURI.RUnlock()
	return calls
}

// GetScope calls GetScopeFunc.
func (mock *AuthorizationCodeMock) GetScope() string {
	if mock.GetScopeFunc == nil {
		panic("AuthorizationCodeMock.GetScopeFunc: method is nil but AuthorizationCode.GetScope was just called")
	}
	callInfo := struct {
	}{}
	mock.lockGetScope.Lock()
	mock.calls.GetScope = append(mock.calls.GetScope, callInfo)
	mock.lockGetScope.Unlock()
	return mock.GetScopeFunc()
}

// GetScopeCalls gets all the calls that were made to GetScope.
// Check the length with:
//
//	len(mockedAuthorizationCode.GetScopeCalls())
func (mock *AuthorizationCodeMock) GetScopeCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockGetScope.RLock()
	calls = mock.calls.GetScope
	mock.lockGetScope.RUnlock()
	return calls
}

// GetUserID calls GetUserIDFunc.
func (mock *AuthorizationCodeMock) GetUserID() uuid.UUID {
	if mock.GetUserIDFunc == nil {
		panic("AuthorizationCodeMock.GetUserIDFunc: method is nil but AuthorizationCode.GetUserID was just called")
	}
	callInfo := struct {
	}{}
	mock.lockGetUserID.Lock()
	mock.calls.GetUserID = append(mock.calls.GetUserID, callInfo)
	mock.lockGetUserID.Unlock()
	return mock.GetUserIDFunc()
}

// GetUserIDCalls gets all the calls that were made to GetUserID.
// Check the length with:
//
//	len(mockedAuthorizationCode.GetUserIDCalls())
func (mock *AuthorizationCodeMock) GetUserIDCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockGetUserID.RLock()
	calls = mock.calls.GetUserID
	mock.lockGetUserID.RUnlock()
	return calls
}

// IsExpired calls IsExpiredFunc.
func (mock *AuthorizationCodeMock) IsExpired(t time.Time) bool {
	if mock.IsExpiredFunc == nil {
		panic("AuthorizationCodeMock.IsExpiredFunc: method is nil but AuthorizationCode.IsExpired was just called")
	}
	callInfo := struct {
		T time.Time
	}{
		T: t,
	}
	mock.lockIsExpired.Lock()
	mock.calls.IsExpired = append(mock.calls.IsExpired, callInfo)
	mock.lockIsExpired.Unlock()
	return mock.IsExpiredFunc(t)
}

// IsExpiredCalls gets all the calls that were made to IsExpired.
// Check the length with:
//
//	len(mockedAuthorizationCode.IsExpiredCalls())
func (mock *AuthorizationCodeMock) IsExpiredCalls() []struct {
	T time.Time
} {
	var calls []struct {
		T time.Time
	}
	mock.lockIsExpired.RLock()
	calls = mock.calls.IsExpired
	mock.lockIsExpired.RUnlock()
	return calls
}

// IsNotFound calls IsNotFoundFunc.
func (mock *AuthorizationCodeMock) IsNotFound() bool {
	if mock.IsNotFoundFunc == nil {
		panic("AuthorizationCodeMock.IsNotFoundFunc: method is nil but AuthorizationCode.IsNotFound was just called")
	}
	callInfo := struct {
	}{}
	mock.lockIsNotFound.Lock()
	mock.calls.IsNotFound = append(mock.calls.IsNotFound, callInfo)
	mock.lockIsNotFound.Unlock()
	return mock.IsNotFoundFunc()
}

// IsNotFoundCalls gets all the calls that were made to IsNotFound.
// Check the length with:
//
//	len(mockedAuthorizationCode.IsNotFoundCalls())
func (mock *AuthorizationCodeMock) IsNotFoundCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockIsNotFound.RLock()
	calls = mock.calls.IsNotFound
	mock.lockIsNotFound.RUnlock()
	return calls
}
