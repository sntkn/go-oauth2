package session

import (
	"encoding/json"
	"fmt"
	"time"

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
func (s *Session) GetSessionData(c *gin.Context, key string, d any) error {
	fullKey := fmt.Sprintf("%s:%s", s.SessionID, key)
	b, err := s.SessionStore.Get(c, fullKey).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(b, &d)
}

// セッションデータをRedisに書き込む関数
func (s *Session) SetSessionData(c *gin.Context, key string, input any) error {
	d, err := json.Marshal(input)
	if err != nil {
		return err
	}
	fullKey := fmt.Sprintf("%s:%s", s.SessionID, key)
	// Redisにセッションデータを書き込み
	return s.SessionStore.Set(c, fullKey, d, 0).Err()
}
