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

func TestLoadSavePopTypedSessionData(t *testing.T) {
	type payload struct {
		Name string
		Age  int
	}

	t.Run("Load returns typed data", func(t *testing.T) {
		input := payload{Name: "john", Age: 10}
		encoded, err := json.Marshal(input)
		require.NoError(t, err)

		mockValkey := &valkey.ClientIFMock{
			GetFunc: func(_ context.Context, _ string) (string, error) {
				return string(encoded), nil
			},
		}

		sess := &Session{SessionID: "sid", SessionStore: mockValkey}
		c := setupTestContext()

		actual, ok, err := Load[payload](c, sess, "typed")
		require.NoError(t, err)
		assert.True(t, ok)
		assert.Equal(t, input, actual)
	})

	t.Run("Save marshals typed payload", func(t *testing.T) {
		input := payload{Name: "mike", Age: 20}
		mockValkey := &valkey.ClientIFMock{
			SetFunc: func(_ context.Context, _ string, value string, _ int64) error {
				var decoded payload
				require.NoError(t, json.Unmarshal([]byte(value), &decoded))
				assert.Equal(t, input, decoded)
				return nil
			},
		}

		sess := &Session{SessionID: "sid", SessionStore: mockValkey}
		c := setupTestContext()

		require.NoError(t, Save(c, sess, "typed", input))
	})

	t.Run("Pop loads and deletes", func(t *testing.T) {
		input := payload{Name: "doe", Age: 30}
		encoded, err := json.Marshal(input)
		require.NoError(t, err)

		mockValkey := &valkey.ClientIFMock{
			GetFunc: func(_ context.Context, _ string) (string, error) {
				return string(encoded), nil
			},
			DelFunc: func(_ context.Context, _ string) error {
				return nil
			},
		}

		sess := &Session{SessionID: "sid", SessionStore: mockValkey}
		c := setupTestContext()

		actual, ok, err := Pop[payload](c, sess, "typed")
		require.NoError(t, err)
		assert.True(t, ok)
		assert.Equal(t, input, actual)
	})
}
