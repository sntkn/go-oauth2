package session

import (
	"time"

	"github.com/cockroachdb/errors"
	"github.com/gin-gonic/gin"
	"github.com/sntkn/go-oauth2/oauth2/internal/redis"
)

type Session struct {
	SessionID    string
	SessionStore *redis.RedisCli
}

func NewSession(c *gin.Context, r *redis.RedisCli) *Session {
	// セッションIDをクッキーから取得
	sessionID, err := c.Cookie("sessionID")
	if err != nil {
		// セッションIDがない場合は新しいセッションIDを生成
		sessionID = GenerateSessionID()
		// クッキーにセッションIDをセット
		c.SetCookie("sessionID", sessionID, 3600, "/", "localhost", false, true)
	}

	// Redisからセッションデータを取得
	sessionData, err := r.Get(c, sessionID).Result()
	if err != nil {
		// セッションデータが存在しない場合は空のデータをセット
		sessionData = ""
	}

	// セッションデータをコンテキストにセット
	c.Set("sessionData", sessionData)

	return &Session{
		SessionID:    sessionID,
		SessionStore: r,
	}
}

// セッションIDを生成する関数
func GenerateSessionID() string {
	return time.Now().Format("20060102150405")
}

// セッションデータを取得する関数
func (s *Session) GetSessionData(c *gin.Context) ([]byte, error) {
	b, err := s.SessionStore.Get(c, s.SessionID).Bytes()
	return b, errors.WithStack(err)
}

// セッションデータをRedisに書き込む関数
func (s *Session) SetSessionData(c *gin.Context, sessionData any) error {
	// Redisにセッションデータを書き込み
	return errors.WithStack(s.SessionStore.Set(c, s.SessionID, sessionData, 0).Err())
}
