package authorization

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sntkn/go-oauth2/oauth2/internal/redis"
	"github.com/sntkn/go-oauth2/oauth2/internal/repository"
	"github.com/sntkn/go-oauth2/oauth2/internal/session"
	"golang.org/x/crypto/bcrypt"
)

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
}

func NewUseCase(redisCli *redis.RedisCli, db *repository.Repository) *UseCase {
	return &UseCase{
		redisCli: redisCli,
		db:       db,
	}
}

func (u *UseCase) Run(c *gin.Context) {
	s := session.NewSession(c, u.redisCli)
	var input AuthorizationInput
	// リクエストのJSONデータをAuthorizationInputにバインド
	if err := c.Bind(&input); err != nil {
		err := fmt.Errorf("Could not bind JSON")
		c.Error(err)
		c.HTML(http.StatusBadRequest, "400.html", gin.H{"error": err})
		return
	}

	if input.Email == "" {
		// TODO: redirect to autorize with parameters
		err := fmt.Errorf("Invalid email address")
		c.Error(err)
		c.HTML(http.StatusBadRequest, "400.html", gin.H{"error": err})
		return
	}

	if input.Password == "" {
		// TODO: redirect to autorize with parameters
		err := fmt.Errorf("Invalid password")
		c.Error(err)
		c.HTML(http.StatusBadRequest, "400.html", gin.H{"error": err})
		return
	}

	// validate user credentials
	user, err := u.db.FindUserByEmail(input.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			// TODO: redirect to autorize with parameters
			c.HTML(http.StatusBadRequest, "400.html", gin.H{"error": err})
		} else {
			c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": err})
		}
		return
	}

	// パスワードを比較して認証
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		// TODO: redirect to autorize with parameters
		c.HTML(http.StatusBadRequest, "400.html", gin.H{"error": err})
		return
	}

	sessionData, err := s.GetSessionData(c)
	if err != nil {
		c.HTML(http.StatusBadRequest, "400.html", gin.H{"error": err})
		return
	}

	var d AuthorizeInput
	err = json.Unmarshal(sessionData, &d)
	if err != nil {
		c.HTML(http.StatusBadRequest, "400.html", gin.H{"error": err})
		return
	}

	// create code
	expired := time.Now().AddDate(0, 0, 10)
	randomString, err := generateRandomString(32)
	if err != nil {
		c.HTML(http.StatusBadRequest, "400.html", gin.H{"error": err})
		return
	}

	clientID, err := uuid.Parse(d.ClientID)
	if err != nil {
		c.HTML(http.StatusBadRequest, "400.html", gin.H{"error": err})
		return
	}

	err = u.db.RegisterOAuth2Code(repository.Code{
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
		c.HTML(http.StatusBadRequest, "400.html", gin.H{"error": err})
		return
	}

	// TODO: clear session data

	c.Redirect(http.StatusFound, fmt.Sprintf("%s?code=%s", d.RedirectURI, randomString))
}

func generateRandomString(length int) (string, error) {
	// ランダムなバイト列を生成
	randomBytes := make([]byte, length)
	_, err := io.ReadFull(rand.Reader, randomBytes)
	if err != nil {
		return "", err
	}

	// URLセーフなBase64エンコード
	encodedString := base64.URLEncoding.EncodeToString(randomBytes)

	return encodedString, nil
}
