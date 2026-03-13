package session

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sntkn/go-oauth2/oauth2/pkg/errors"
	"github.com/sntkn/go-oauth2/oauth2/pkg/valkey"
)

type Creator func(c *gin.Context) *Session

const sessionExpirationSeconds = 3600

//go:generate go run github.com/matryer/moq -out session_manager_mock.go . SessionManager
type SessionManager interface {
	NewSession(c *gin.Context) SessionClient
}

type DefaultSessionManager struct {
	cli     valkey.ClientIF
	expires int
}

func NewSessionManager(cli valkey.ClientIF, expires int) *DefaultSessionManager {
	return &DefaultSessionManager{
		cli:     cli,
		expires: expires,
	}
}

func (m *DefaultSessionManager) NewSession(c *gin.Context) SessionClient {
	// セッションIDをクッキーから取得
	sessionID, err := c.Cookie("sessionID")
	if err != nil {
		// セッションIDがない場合は新しいセッションIDを生成
		sessionID = GenerateSessionID()
		// クッキーにセッションIDをセット
		c.SetCookie("sessionID", sessionID, m.expires, "/", "localhost", false, true)
	}

	return &Session{
		SessionID:    sessionID,
		SessionStore: m.cli,
	}
}

type Session struct {
	SessionID    string
	SessionStore valkey.ClientIF
}

//go:generate go run github.com/matryer/moq -out session_mock.go . SessionClient
type SessionClient interface {
	GetSessionData(c *gin.Context, key string) (string, error)
	SetSessionData(c *gin.Context, key string, input string) error
	DelSessionData(c *gin.Context, key string) error
	PullSessionData(c *gin.Context, key string) (string, error)
}

// Load retrieves typed session data.
func Load[T any](c *gin.Context, s SessionClient, key string) (T, bool, error) {
	var zero T
	data, err := s.GetSessionData(c, key)
	if err != nil {
		return zero, false, err
	}
	if len(data) == 0 {
		return zero, false, nil
	}
	var value T
	if err := json.Unmarshal([]byte(data), &value); err != nil {
		return zero, false, errors.WithStack(err)
	}
	return value, true, nil
}

// Save stores typed session data.
func Save[T any](c *gin.Context, s SessionClient, key string, value T) error {
	payload, err := json.Marshal(value)
	if err != nil {
		return errors.WithStack(err)
	}
	return s.SetSessionData(c, key, string(payload))
}

// Pop loads and deletes typed session data in a single call.
func Pop[T any](c *gin.Context, s SessionClient, key string) (T, bool, error) {
	value, ok, err := Load[T](c, s, key)
	if err != nil || !ok {
		return value, ok, err
	}
	if err := s.DelSessionData(c, key); err != nil {
		var zero T
		return zero, false, err
	}
	return value, true, nil
}

// セッションIDを生成する関数
func GenerateSessionID() string {
	return time.Now().Format("20060102150405")
}

// セッションデータを取得する関数
func (s *Session) GetSessionData(c *gin.Context, key string) (string, error) {
	return s.SessionStore.Get(c, s.fullKey(key))
}

func (s *Session) SetSessionData(c *gin.Context, key string, input string) error {
	return s.SessionStore.Set(c, s.fullKey(key), input, sessionExpirationSeconds)
}

func (s *Session) DelSessionData(c *gin.Context, key string) error {
	return s.SessionStore.Del(c, s.fullKey(key))
}

func (s *Session) fullKey(key string) string {
	return fmt.Sprintf("%s:%s", s.SessionID, key)
}

// Get and flush session
func (s *Session) PullSessionData(c *gin.Context, key string) (string, error) {
	v, err := s.GetSessionData(c, key)
	if err != nil {
		return "", err
	}

	if err := s.DelSessionData(c, key); err != nil {
		return "", err
	}

	return v, nil
}

// func GetSessionDataToType[T any](s *Session, c *gin.Context, key string, t T) (T, error) {
//	b, err := s.GetSessionData(c, key)
//	if err != nil {
//		return t, err
//	}
//	err = json.Unmarshal(b, t)
//	return t, err
//}
