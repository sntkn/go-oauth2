package usecases

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sntkn/go-oauth2/oauth2/internal/repository"
	"github.com/sntkn/go-oauth2/oauth2/internal/session"
	"github.com/sntkn/go-oauth2/oauth2/pkg/config"
	cerrs "github.com/sntkn/go-oauth2/oauth2/pkg/errors"
	"github.com/sntkn/go-oauth2/oauth2/pkg/redis"
	"golang.org/x/crypto/bcrypt"
)

type Authorization struct {
	redisCli *redis.RedisCli
	db       *repository.Repository
	cfg      *config.Config
	sess     *session.Session
}

func NewAuthorization(redisCli *redis.RedisCli,
	db *repository.Repository,
	cfg *config.Config,
	sess *session.Session,
) *Authorization {
	return &Authorization{
		redisCli: redisCli,
		db:       db,
		cfg:      cfg,
		sess:     sess,
	}
}

func (u *Authorization) Invoke(c *gin.Context, email string, password string) (string, error) {

	// validate user credentials
	user, err := u.db.FindUserByEmail(email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", cerrs.NewUsecaseError(http.StatusFound, err.Error())
		}
		return "", cerrs.NewUsecaseError(http.StatusInternalServerError, err.Error())
	}

	// パスワードを比較して認証
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", cerrs.NewUsecaseError(http.StatusFound, err.Error())
	}
	var d AuthorizeInput
	if err := u.sess.GetNamedSessionData(c, "auth", &d); err != nil {
		return "", cerrs.NewUsecaseError(http.StatusInternalServerError, err.Error())
	}

	// create code
	expired := time.Now().Add(u.cfg.AuthCodeExpires * time.Second)
	randomStringLen := 32
	randomString, err := generateRandomString(randomStringLen)
	if err != nil {
		return "", cerrs.NewUsecaseError(http.StatusInternalServerError, err.Error())
	}

	clientID, err := uuid.Parse(d.ClientID)
	if err != nil {
		return "", cerrs.NewUsecaseError(http.StatusInternalServerError, err.Error())
	}

	err = u.db.RegisterOAuth2Code(&repository.Code{
		Code:        randomString,
		ClientID:    clientID,
		UserID:      user.ID,
		Scope:       d.Scope,
		RedirectURI: d.RedirectURI,
		ExpiresAt:   expired,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	})
	if err != nil {
		return "", cerrs.NewUsecaseError(http.StatusInternalServerError, err.Error())
	}

	if err := u.sess.DelSessionData(c, "auth"); err != nil {
		return "", cerrs.NewUsecaseError(http.StatusInternalServerError, err.Error())
	}

	return fmt.Sprintf("%s?code=%s", d.RedirectURI, randomString), nil
}

func generateRandomString(length int) (string, error) {
	// ランダムなバイト列を生成
	randomBytes := make([]byte, length)
	_, err := io.ReadFull(rand.Reader, randomBytes)
	if err != nil {
		return "", errors.WithStack(err)
	}

	// URLセーフなBase64エンコード
	encodedString := base64.URLEncoding.EncodeToString(randomBytes)

	return encodedString, nil
}
