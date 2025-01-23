package session

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sntkn/go-oauth2/oauth2/pkg/valkey"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestContext() *gin.Context {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(nil)
	return c
}

func TestGetSessionData(t *testing.T) {
	mockValkey := &valkey.ClientIFMock{
		GetFunc: func(_ context.Context, _ string) (string, error) {
			return "testValue", nil
		},
	}

	c := setupTestContext()

	session := &Session{
		SessionID:    "sessionID",
		SessionStore: mockValkey,
	}

	// Test GetSessionData
	value, err := session.GetSessionData(c, "testKey")
	require.NoError(t, err)
	assert.Equal(t, "testValue", value)
}

func TestSetSessionData(t *testing.T) {
	mockValkey := &valkey.ClientIFMock{
		SetFunc: func(_ context.Context, _ string, _ string, _ int64) error {
			return nil
		},
	}

	c := setupTestContext()

	session := &Session{
		SessionID:    "sessionID",
		SessionStore: mockValkey,
	}

	// Test GetSessionData
	err := session.SetSessionData(c, "testKey", "testValue")
	require.NoError(t, err)
}

func TestDelSessionData(t *testing.T) {
	mockValkey := &valkey.ClientIFMock{
		DelFunc: func(_ context.Context, _ string) error {
			return nil
		},
	}

	c := setupTestContext()

	session := &Session{
		SessionID:    "sessionID",
		SessionStore: mockValkey,
	}

	// Test GetSessionData
	err := session.DelSessionData(c, "testKey")
	require.NoError(t, err)
}

func TestPullSessionData(t *testing.T) {
	mockValkey := &valkey.ClientIFMock{
		GetFunc: func(_ context.Context, _ string) (string, error) {
			return "testValue", nil
		},
		DelFunc: func(_ context.Context, _ string) error {
			return nil
		},
	}

	c := setupTestContext()

	session := &Session{
		SessionID:    "sessionID",
		SessionStore: mockValkey,
	}

	// Test GetSessionData
	value, err := session.PullSessionData(c, "testKey")
	require.NoError(t, err)
	assert.Equal(t, "testValue", value)
}

func TestGetNamedSessionData(t *testing.T) {
	type TestStruct struct {
		Name  string
		Value int
	}

	input := TestStruct{
		Name:  "test",
		Value: 42,
	}

	data, err := json.Marshal(input)
	require.NoError(t, err)

	mockValkey := &valkey.ClientIFMock{
		GetFunc: func(_ context.Context, _ string) (string, error) {
			return string(data), nil
		},
	}

	c := setupTestContext()

	session := &Session{
		SessionID:    "sessionID",
		SessionStore: mockValkey,
	}

	var output TestStruct

	// Test GetSessionData
	err = session.GetNamedSessionData(c, "testKey", &output)
	require.NoError(t, err)
	assert.Equal(t, input, output)
}

func TestSetNamedSessionData(t *testing.T) {
	type TestStruct struct {
		Name  string
		Value int
	}

	input := TestStruct{
		Name:  "test",
		Value: 42,
	}

	data, err := json.Marshal(input)
	require.NoError(t, err)

	mockValkey := &valkey.ClientIFMock{
		SetFunc: func(_ context.Context, _ string, _ string, _ int64) error {
			return nil
		},
	}

	c := setupTestContext()

	session := &Session{
		SessionID:    "sessionID",
		SessionStore: mockValkey,
	}

	// Test SetSessionData
	err = session.SetNamedSessionData(c, "testKey", data)
	require.NoError(t, err)
}

func TestFlushNamedSessionData(t *testing.T) {
	t.Parallel()

	type TestStruct struct {
		Name  string
		Value int
	}

	input := TestStruct{
		Name:  "test",
		Value: 42,
	}

	data, err := json.Marshal(input)
	require.NoError(t, err)

	mockValkey := &valkey.ClientIFMock{
		GetFunc: func(_ context.Context, _ string) (string, error) {
			return string(data), nil
		},
		DelFunc: func(_ context.Context, _ string) error {
			return nil
		},
	}

	c := setupTestContext()

	session := &Session{
		SessionID:    "sessionID",
		SessionStore: mockValkey,
	}

	var output TestStruct

	err = session.FlushNamedSessionData(c, "testKey", &output)
	require.NoError(t, err)
	assert.Equal(t, input, output)
}
