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
	GetNamedSessionData(c *gin.Context, key string, t any) error
	SetNamedSessionData(c *gin.Context, key string, v any) error
	FlushNamedSessionData(c *gin.Context, key string, t any) error
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
	return s.SessionStore.Set(c, s.fullKey(key), input, 3600)
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

func (s *Session) GetNamedSessionData(c *gin.Context, key string, t any) error {
	str, err := s.GetSessionData(c, key)
	if err != nil {
		return err
	}
	if len(str) == 0 {
		return nil
	}
	if err = json.Unmarshal([]byte(str), &t); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (s *Session) SetNamedSessionData(c *gin.Context, key string, v any) error {
	d, err := json.Marshal(v)
	if err != nil {
		return errors.WithStack(err)
	}
	return s.SetSessionData(c, key, string(d))
}

func (s *Session) FlushNamedSessionData(c *gin.Context, key string, t any) error {
	if err := s.GetNamedSessionData(c, key, t); err != nil {
		return err
	}
	if err := s.DelSessionData(c, key); err != nil {
		return err
	}
	return nil
}

// func GetSessionDataToType[T any](s *Session, c *gin.Context, key string, t T) (T, error) {
//	b, err := s.GetSessionData(c, key)
//	if err != nil {
//		return t, err
//	}
//	err = json.Unmarshal(b, t)
//	return t, err
//}
