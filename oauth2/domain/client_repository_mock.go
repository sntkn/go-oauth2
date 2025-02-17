// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package domain

import (
	"github.com/google/uuid"
	"sync"
)

// Ensure, that ClientRepositoryMock does implement ClientRepository.
// If this is not the case, regenerate this file with moq.
var _ ClientRepository = &ClientRepositoryMock{}

// ClientRepositoryMock is a mock implementation of ClientRepository.
//
//	func TestSomethingThatUsesClientRepository(t *testing.T) {
//
//		// make and configure a mocked ClientRepository
//		mockedClientRepository := &ClientRepositoryMock{
//			FindClientByClientIDFunc: func(clientID uuid.UUID) (Client, error) {
//				panic("mock out the FindClientByClientID method")
//			},
//		}
//
//		// use mockedClientRepository in code that requires ClientRepository
//		// and then make assertions.
//
//	}
type ClientRepositoryMock struct {
	// FindClientByClientIDFunc mocks the FindClientByClientID method.
	FindClientByClientIDFunc func(clientID uuid.UUID) (Client, error)

	// calls tracks calls to the methods.
	calls struct {
		// FindClientByClientID holds details about calls to the FindClientByClientID method.
		FindClientByClientID []struct {
			// ClientID is the clientID argument value.
			ClientID uuid.UUID
		}
	}
	lockFindClientByClientID sync.RWMutex
}

// FindClientByClientID calls FindClientByClientIDFunc.
func (mock *ClientRepositoryMock) FindClientByClientID(clientID uuid.UUID) (Client, error) {
	if mock.FindClientByClientIDFunc == nil {
		panic("ClientRepositoryMock.FindClientByClientIDFunc: method is nil but ClientRepository.FindClientByClientID was just called")
	}
	callInfo := struct {
		ClientID uuid.UUID
	}{
		ClientID: clientID,
	}
	mock.lockFindClientByClientID.Lock()
	mock.calls.FindClientByClientID = append(mock.calls.FindClientByClientID, callInfo)
	mock.lockFindClientByClientID.Unlock()
	return mock.FindClientByClientIDFunc(clientID)
}

// FindClientByClientIDCalls gets all the calls that were made to FindClientByClientID.
// Check the length with:
//
//	len(mockedClientRepository.FindClientByClientIDCalls())
func (mock *ClientRepositoryMock) FindClientByClientIDCalls() []struct {
	ClientID uuid.UUID
} {
	var calls []struct {
		ClientID uuid.UUID
	}
	mock.lockFindClientByClientID.RLock()
	calls = mock.calls.FindClientByClientID
	mock.lockFindClientByClientID.RUnlock()
	return calls
}
