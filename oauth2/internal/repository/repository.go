package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type Repository struct {
	db *sql.DB
}

type Conn struct {
	Host     string
	Port     uint32
	User     string
	Password string
	DBName   string
}

type User struct {
	ID       uuid.UUID `db:"id"`
	Email    string    `db:"email"`
	Password string    `db:"password"`
	// 他のユーザー属性をここに追加
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
	ExpiresAt   time.Time `db:"expired_at"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type Token struct {
	AccessToken string    `db:"access_token"`
	ClientID    uuid.UUID `db:"client_id"`
	UserID      uuid.UUID `db:"user_id"`
	Scope       string    `db:"scope"`
	ExpiresAt   time.Time `db:"expired_at"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type RefreshToken struct {
	RefreshToken string    `db:"refresh_token"`
	AccessToken  string    `db:"access_token"`
	ExpiresAt    time.Time `db:"expired_at"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

func NewClient(c Conn) (*Repository, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.Host, c.Port, c.User, c.Password, c.DBName)

	// PostgreSQLに接続
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	return &Repository{
		db: db,
	}, nil
}

func (c *Repository) Close() {
	c.db.Close()
}

func (r *Repository) FindClientByClientID(clientID string) (Client, error) {
	q := "SELECT id, redirect_uris FROM oauth2_clients WHERE id = $1"
	var client Client

	err := r.db.QueryRow(q, clientID).Scan(&client.ID, &client.RedirectURIs)
	return client, err
}

func (r *Repository) FindUserByEmail(email string) (User, error) {
	q := "SELECT id, email, password FROM users WHERE email = $1"
	var user User

	err := r.db.QueryRow(q, email).Scan(&user.ID, &user.Email, &user.Password)
	return user, err
}

func (r *Repository) RegisterOAuth2Code(c Code) error {
	q := `
			INSERT INTO oauth2_codes
				(code, client_id, user_id, scope, redirect_uri, expires_at, created_at, updated_at)
			VALUES
				($1, $2, $3, $4, $5, $6, $7, $8)
		`
	_, err := r.db.Exec(q, c.Code, c.ClientID, c.UserID, c.Scope, c.RedirectURI, c.ExpiresAt, time.Now(), time.Now())
	return err
}

func (r *Repository) FindValidOAuth2Code(code string, expiresAt time.Time) (Code, error) {
	q := "SELECT user_id, client_id, scope, expires_at FROM oauth2_codes WHERE code = $1 AND revoked_at IS NULL AND expires_at > $2"
	var c Code

	err := r.db.QueryRow(q, code, expiresAt).Scan(&c.UserID, &c.ClientID, &c.Scope, &c.ExpiresAt)
	return c, err
}

func (r *Repository) RegisterToken(t Token) error {
	q := "INSERT INTO oauth2_tokens (access_token, client_id, user_id, scope, expires_at, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	_, err := r.db.Exec(q, t.AccessToken, t.ClientID, t.UserID, t.Scope, t.ExpiresAt, time.Now(), time.Now())
	return err
}

func (r *Repository) RegesterRefreshToken(t RefreshToken) error {
	q := "INSERT INTO oauth2_refresh_tokens (refresh_token, access_token, expires_at, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)"
	_, err := r.db.Exec(q, t.RefreshToken, t.AccessToken, t.ExpiresAt, time.Now(), time.Now())
	return err
}

func (r *Repository) RevokeCode(code string) error {
	updateQuery := "UPDATE oauth2_codes SET revoked_at = $1 WHERE code = $2"
	_, err := r.db.Exec(updateQuery, time.Now(), code)
	return err
}

func (r *Repository) FindValidRefreshToken(refreshToken string, expiresAt time.Time) (RefreshToken, error) {
	q := "SELECT refresh_token FROM oauth2_refresh_tokens WHERE refresh_token = $1 AND expires_at > $2"
	var rtkn RefreshToken

	err := r.db.QueryRow(q, refreshToken, expiresAt).Scan(&rtkn.RefreshToken)
	return rtkn, err
}
