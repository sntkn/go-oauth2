package session

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sntkn/go-oauth2/oauth2/pkg/errors"
	"github.com/sntkn/go-oauth2/oauth2/pkg/redis"
)

type Creator func(c *gin.Context) *Session

type SessionManager interface {
	NewSession(c *gin.Context) *Session
}

type DefaultSessionManager struct {
	cli     redis.RedisClient
	expires int
}

func NewSessionManager(cli redis.RedisClient, expires int) *DefaultSessionManager {
	return &DefaultSessionManager{
		cli:     cli,
		expires: expires,
	}
}

func (m *DefaultSessionManager) NewSession(c *gin.Context) *Session {
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
	SessionStore redis.RedisClient
}

//go:generate go run github.com/matryer/moq -out session_mock.go . SessionClient
type SessionClient interface {
	GetSessionData(c *gin.Context, key string) ([]byte, error)
	SetSessionData(c *gin.Context, key string, input any) error
	DelSessionData(c *gin.Context, key string) error
	PullSessionData(c *gin.Context, key string) ([]byte, error)
	GetNamedSessionData(c *gin.Context, key string, t any) error
	SetNamedSessionData(c *gin.Context, key string, v any) error
	FlushNamedSessionData(c *gin.Context, key string, t any) error
}

// セッションIDを生成する関数
func GenerateSessionID() string {
	return time.Now().Format("20060102150405")
}

// セッションデータを取得する関数
func (s *Session) GetSessionData(c *gin.Context, key string) ([]byte, error) {
	fullKey := fmt.Sprintf("%s:%s", s.SessionID, key)
	b, err := s.SessionStore.GetOrNil(c, fullKey)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return b, nil
}

// セッションデータをRedisに書き込む関数
func (s *Session) SetSessionData(c *gin.Context, key string, input any) error {
	fullKey := fmt.Sprintf("%s:%s", s.SessionID, key)
	// Redisにセッションデータを書き込み
	return s.SessionStore.Set(c, fullKey, input, 0)
}

func (s *Session) DelSessionData(c *gin.Context, key string) error {
	fullKey := fmt.Sprintf("%s:%s", s.SessionID, key)
	return s.SessionStore.Del(c, fullKey)
}

// Get and flush session
func (s *Session) PullSessionData(c *gin.Context, key string) ([]byte, error) {
	v, err := s.GetSessionData(c, key)
	if err != nil {
		return nil, err
	}

	if err := s.DelSessionData(c, key); err != nil {
		return nil, err
	}

	return v, nil
}

func (s *Session) GetNamedSessionData(c *gin.Context, key string, t any) error {
	b, err := s.GetSessionData(c, key)
	if err != nil {
		return errors.WithStack(err)
	}
	if len(b) == 0 {
		return nil
	}
	if err = json.Unmarshal(b, &t); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (s *Session) SetNamedSessionData(c *gin.Context, key string, v any) error {
	d, err := json.Marshal(v)
	if err != nil {
		return errors.WithStack(err)
	}
	return s.SetSessionData(c, key, d)
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
