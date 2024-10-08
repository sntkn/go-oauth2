// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package auth

import (
	"sync"
)

// Ensure, that AuthorizeUsecaseMock does implement AuthorizeUsecase.
// If this is not the case, regenerate this file with moq.
var _ AuthorizeUsecase = &AuthorizeUsecaseMock{}

// AuthorizeUsecaseMock is a mock implementation of AuthorizeUsecase.
//
//	func TestSomethingThatUsesAuthorizeUsecase(t *testing.T) {
//
//		// make and configure a mocked AuthorizeUsecase
//		mockedAuthorizeUsecase := &AuthorizeUsecaseMock{
//			InvokeFunc: func(clientID string, redirectURI string) error {
//				panic("mock out the Invoke method")
//			},
//		}
//
//		// use mockedAuthorizeUsecase in code that requires AuthorizeUsecase
//		// and then make assertions.
//
//	}
type AuthorizeUsecaseMock struct {
	// InvokeFunc mocks the Invoke method.
	InvokeFunc func(clientID string, redirectURI string) error

	// calls tracks calls to the methods.
	calls struct {
		// Invoke holds details about calls to the Invoke method.
		Invoke []struct {
			// ClientID is the clientID argument value.
			ClientID string
			// RedirectURI is the redirectURI argument value.
			RedirectURI string
		}
	}
	lockInvoke sync.RWMutex
}

// Invoke calls InvokeFunc.
func (mock *AuthorizeUsecaseMock) Invoke(clientID string, redirectURI string) error {
	if mock.InvokeFunc == nil {
		panic("AuthorizeUsecaseMock.InvokeFunc: method is nil but AuthorizeUsecase.Invoke was just called")
	}
	callInfo := struct {
		ClientID    string
		RedirectURI string
	}{
		ClientID:    clientID,
		RedirectURI: redirectURI,
	}
	mock.lockInvoke.Lock()
	mock.calls.Invoke = append(mock.calls.Invoke, callInfo)
	mock.lockInvoke.Unlock()
	return mock.InvokeFunc(clientID, redirectURI)
}

// InvokeCalls gets all the calls that were made to Invoke.
// Check the length with:
//
//	len(mockedAuthorizeUsecase.InvokeCalls())
func (mock *AuthorizeUsecaseMock) InvokeCalls() []struct {
	ClientID    string
	RedirectURI string
} {
	var calls []struct {
		ClientID    string
		RedirectURI string
	}
	mock.lockInvoke.RLock()
	calls = mock.calls.Invoke
	mock.lockInvoke.RUnlock()
	return calls
}