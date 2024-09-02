package usecases

import (
	"database/sql"
	"net/http"

	"github.com/cockroachdb/errors"
	"github.com/gin-gonic/gin"
	"github.com/sntkn/go-oauth2/oauth2/internal/repository"
	"github.com/sntkn/go-oauth2/oauth2/pkg/config"
	cerrs "github.com/sntkn/go-oauth2/oauth2/pkg/errors"
)

type AuthorizeInput struct {
	ResponseType string `form:"response_type" binding:"required"`
	ClientID     string `form:"client_id" binding:"required"`
	Scope        string `form:"scope" binding:"required"`
	RedirectURI  string `form:"redirect_uri" binding:"required"`
	State        string `form:"state" binding:"required"`
}

type Authorize struct {
	cfg *config.Config
	db  repository.OAuth2Repository
}

func NewAuthorize(cfg *config.Config, db repository.OAuth2Repository) *Authorize {
	return &Authorize{
		cfg: cfg,
		db:  db,
	}
}

func (u *Authorize) Invoke(_ *gin.Context, clientID, redirectURI string) error {
	// check client
	client, err := u.db.FindClientByClientID(clientID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return cerrs.NewUsecaseError(http.StatusBadRequest, err.Error())
		}
		return cerrs.NewUsecaseError(http.StatusInternalServerError, err.Error())
	}

	if client.RedirectURIs != redirectURI {
		return cerrs.NewUsecaseError(http.StatusBadRequest, "redirect uri does not match")
	}

	return nil
}
