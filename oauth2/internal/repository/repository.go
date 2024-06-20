package repository

import (
	"fmt"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type Repository struct {
	db *sqlx.DB
}

type Conn struct {
	Host     string
	Port     uint32
	User     string
	Password string
	DBName   string
}

type User struct {
	ID        uuid.UUID `db:"id"`
	Name      string    `db:"name"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type Client struct {
	ID           uuid.UUID `db:"id"`
	Name         string    `db:"name"`
	RedirectURIs string    `db:"redirect_uris"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

type Code struct {
	Code        string    `db:"code"`
	ClientID    uuid.UUID `db:"client_id"`
	UserID      uuid.UUID `db:"user_id"`
	Scope       string    `db:"scope"`
	RedirectURI string    `db:"redirect_uri"`
	ExpiresAt   time.Time `db:"expires_at"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type Token struct {
	AccessToken string    `db:"access_token"`
	ClientID    uuid.UUID `db:"client_id"`
	UserID      uuid.UUID `db:"user_id"`
	Scope       string    `db:"scope"`
	ExpiresAt   time.Time `db:"expires_at"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type RefreshToken struct {
	RefreshToken string    `db:"refresh_token"`
	AccessToken  string    `db:"access_token"`
	ExpiresAt    time.Time `db:"expires_at"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

func NewClient(c Conn) (*Repository, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.Host, c.Port, c.User, c.Password, c.DBName)

	// PostgreSQLに接続
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &Repository{
		db: db,
	}, nil
}

func (r *Repository) Close() {
	r.db.Close()
}

func (r *Repository) FindClientByClientID(clientID string) (Client, error) {
	q := "SELECT id, redirect_uris FROM oauth2_clients WHERE id = $1"
	var c Client

	err := r.db.Get(&c, q, clientID)
	return c, errors.WithStack(err)
}

func (r *Repository) FindUserByEmail(email string) (User, error) {
	q := "SELECT id, email, password FROM users WHERE email = $1"
	var u User

	err := r.db.Get(&u, q, email)
	return u, errors.WithStack(err)
}

func (r *Repository) RegisterOAuth2Code(c *Code) error {
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
	q := `
			INSERT INTO oauth2_codes
				(code, client_id, user_id, scope, redirect_uri, expires_at, created_at, updated_at)
			VALUES
				(:code, :client_id, :user_id, :scope, :redirect_uri, :expires_at, :created_at, :updated_at)
		`
	_, err := r.db.NamedExec(q, c)
	return errors.WithStack(err)
}

func (r *Repository) FindValidOAuth2Code(code string, expiresAt time.Time) (Code, error) {
	q := "SELECT user_id, client_id, scope, expires_at FROM oauth2_codes WHERE code = $1 AND revoked_at IS NULL AND expires_at > $2"
	var c Code

	err := r.db.Get(&c, q, code, expiresAt)
	return c, errors.WithStack(err)
}

func (r *Repository) RegisterToken(t *Token) error {
	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()
	q := `
		INSERT INTO oauth2_tokens (access_token, client_id, user_id, scope, expires_at, created_at, updated_at)
		VALUES (:access_token, :client_id, :user_id, :scope, :expires_at, :created_at, :updated_at)
	`
	_, err := r.db.NamedExec(q, t)
	return errors.WithStack(err)
}

func (r *Repository) RegesterRefreshToken(t *RefreshToken) error {
	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()
	q := `INSERT INTO oauth2_refresh_tokens (refresh_token, access_token, expires_at, created_at, updated_at)
	VALUES (:refresh_token, :access_token, :expires_at, :created_at, :updated_at)`
	_, err := r.db.NamedExec(q, t)
	return errors.WithStack(err)
}

func (r *Repository) RevokeCode(code string) error {
	updateQuery := "UPDATE oauth2_codes SET revoked_at = $1 WHERE code = $2"
	_, err := r.db.Exec(updateQuery, time.Now(), code)
	return errors.WithStack(err)
}

func (r *Repository) FindValidRefreshToken(refreshToken string, expiresAt time.Time) (RefreshToken, error) {
	q := "SELECT access_token FROM oauth2_refresh_tokens WHERE refresh_token = $1 AND expires_at > $2"
	var rtkn RefreshToken

	err := r.db.Get(&rtkn, q, refreshToken, expiresAt)
	return rtkn, errors.WithStack(err)
}

func (r *Repository) FindToken(accessToken string) (Token, error) {
	q := "SELECT user_id, client_id, scope FROM oauth2_tokens WHERE access_token = $1"
	var tkn Token

	err := r.db.Get(&tkn, q, accessToken)
	return tkn, errors.WithStack(err)
}

func (r *Repository) RevokeToken(accessToken string) error {
	updateQuery := "UPDATE oauth2_tokens SET revoked_at = $1 WHERE access_token = $2"
	_, err := r.db.Exec(updateQuery, time.Now(), accessToken)
	return errors.WithStack(err)
}

func (r *Repository) RevokeRefreshToken(refreshToken string) error {
	updateQuery := "UPDATE oauth2_refresh_tokens SET revoked_at = $1 WHERE refresh_token = $2"
	_, err := r.db.Exec(updateQuery, time.Now(), refreshToken)
	return errors.WithStack(err)
}

func (r *Repository) FindUser(id uuid.UUID) (User, error) {
	q := "SELECT id, name, email FROM users WHERE id = $1"
	var u User

	err := r.db.Get(&u, q, id)
	return u, errors.WithStack(err)
}

func (r *Repository) ExistsUserByEmail(email string) (bool, error) {
	q := "SELECT count(*) AS c FROM users WHERE email = $1"
	var c uint8

	err := r.db.QueryRow(q, email).Scan(&c)
	return c > 0, errors.WithStack(err)
}

func (r *Repository) CreateUser(u *User) error {
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.WithStack(err)
	}
	u.Password = string(hashedPassword)
	q := `INSERT INTO users (name, email, password, created_at, updated_at)
	VALUES (:name, :email, :password, :created_at, :updated_at)`
	_, err = r.db.NamedExec(q, u)
	return errors.WithStack(err)
}
