package authorize

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/cockroachdb/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sntkn/go-oauth2/oauth2/internal/redis"
	"github.com/sntkn/go-oauth2/oauth2/internal/repository"
	"github.com/sntkn/go-oauth2/oauth2/internal/session"
)

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
	// /authorize?response_type=code&client_id=550e8400-e29b-41d4-a716-446655440000&scope=read&redirect_uri=http%3A%2F%2Flocalhost%3A8000%2Fcallback&state=ok
	var input AuthorizeInput
	// Query ParameterをAuthorizeInputにバインド
	if err := c.BindQuery(&input); err != nil {
		c.Error(errors.WithStack(err))
		c.HTML(http.StatusBadRequest, "400.html", gin.H{"error": err.Error()})
		return
	}

	if input.ResponseType == "" {
		err := fmt.Errorf("Invalid response_type")
		c.Error(errors.WithStack(err))
		c.HTML(http.StatusBadRequest, "400.html", gin.H{"error": err.Error()})
		return
	}
	if input.ResponseType != "code" {
		err := fmt.Errorf("Invalid response_type: code must be 'code'")
		c.Error(errors.WithStack(err))
		c.HTML(http.StatusBadRequest, "400.html", gin.H{"error": err.Error()})
	}

	if input.ClientID == "" {
		err := fmt.Errorf("Invalid client_id")
		c.Error(errors.WithStack(err))
		c.HTML(http.StatusBadRequest, "400.html", gin.H{"error": err.Error()})
		return
	}
	if IsValidUUID(input.ClientID) == false {
		err := fmt.Errorf("Could not parse client_id")
		c.Error(errors.WithStack(err))
		c.HTML(http.StatusBadRequest, "400.html", gin.H{"error": err.Error()})
		return
	}

	if input.Scope == "" {
		err := fmt.Errorf("Invalid scope")
		c.Error(errors.WithStack(err))
		c.HTML(http.StatusBadRequest, "400.html", gin.H{"error": err.Error()})
		return
	}

	if input.RedirectURI == "" {
		err := fmt.Errorf("Invalid redirect_uri")
		c.Error(errors.WithStack(err))
		c.HTML(http.StatusBadRequest, "400.html", gin.H{"error": err.Error()})
		return
	}

	if input.State == "" {
		err := fmt.Errorf("Invalid state")
		c.Error(errors.WithStack(err))
		c.HTML(http.StatusBadRequest, "400.html", gin.H{"error": err.Error()})
		return
	}

	// check client
	client, err := u.db.FindClientByClientID(input.ClientID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.Error(err)
			c.HTML(http.StatusBadRequest, "400.html", gin.H{"error": err.Error()})
		} else {
			c.Error(err)
			c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": err.Error()})
		}
		return
	}

	if client.RedirectURIs != input.RedirectURI {
		err := fmt.Errorf("Redirect URI does not match")
		c.Error(errors.WithStack(err))
		c.HTML(http.StatusBadRequest, "400.html", gin.H{"error": err.Error()})
		return
	}

	// セッションデータを書き込む
	if err = s.SetNamedSessionData(c, "auth", &input); err != nil {
		c.Error(err)
		c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusFound, "/signin")
}

func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}
