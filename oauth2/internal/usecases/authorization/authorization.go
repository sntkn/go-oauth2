package authorization

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
	"github.com/sntkn/go-oauth2/oauth2/internal/redis"
	"github.com/sntkn/go-oauth2/oauth2/internal/repository"
	"github.com/sntkn/go-oauth2/oauth2/internal/session"
	"github.com/sntkn/go-oauth2/oauth2/pkg/config"
	"golang.org/x/crypto/bcrypt"
)

type SigninForm struct {
	Email string `form:"email"`
	Error string
}

type AuthorizationInput struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}

type AuthorizeInput struct {
	ResponseType string `form:"response_type"`
	ClientID     string `form:"client_id"`
	Scope        string `form:"scope"`
	RedirectURI  string `form:"redirect_uri"`
	State        string `form:"state"`
}

type UseCase struct {
	redisCli *redis.RedisCli
	db       *repository.Repository
	cfg      *config.Config
}

func NewUseCase(redisCli *redis.RedisCli, db *repository.Repository, cfg *config.Config) *UseCase {
	return &UseCase{
		redisCli: redisCli,
		db:       db,
		cfg:      cfg,
	}
}

func (u *UseCase) Run(c *gin.Context) {
	s := session.NewSession(c, u.redisCli)
	var input AuthorizationInput
	// リクエストのJSONデータをAuthorizationInputにバインド
	if err := c.Bind(&input); err != nil {
		err := fmt.Errorf("could not bind JSON")
		c.Error(errors.WithStack(err))
		c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": err.Error()})
		return
	}

	if err := s.SetNamedSessionData(c, "signin_form", SigninForm{
		Email: input.Email,
	}); err != nil {
		c.Error(errors.WithStack(err))
		c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": err.Error()})
	}

	if input.Email == "" {
		c.Redirect(http.StatusFound, "/signin")
		return
	}

	if input.Password == "" {
		c.Redirect(http.StatusFound, "/signin")
		return
	}

	// validate user credentials
	user, err := u.db.FindUserByEmail(input.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.Redirect(http.StatusFound, "/signin")
		} else {
			c.Error(err)
			c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": err.Error()})
		}
		return
	}

	// パスワードを比較して認証
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		c.Redirect(http.StatusFound, "/signin")
		return
	}

	var d AuthorizeInput
	if err = s.GetNamedSessionData(c, "auth", &d); err != nil {
		c.Error(err)
		c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": err.Error()})
		return
	}

	// create code
	expired := time.Now().Add(u.cfg.AuthCodeExpires * time.Second)
	randomStringLen := 32
	randomString, err := generateRandomString(randomStringLen)
	if err != nil {
		c.Error(errors.WithStack(err))
		c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": err.Error()})
		return
	}

	clientID, err := uuid.Parse(d.ClientID)
	if err != nil {
		c.Error(errors.WithStack(err))
		c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": err.Error()})
		return
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
		c.Error(err)
		c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": err.Error()})
		return
	}

	if err := s.DelSessionData(c, "auth"); err != nil {
		c.Error(err)
		c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusFound, fmt.Sprintf("%s?code=%s", d.RedirectURI, randomString))
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
