// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameOauth2Token = "oauth2_tokens"

// Oauth2Token mapped from table <oauth2_tokens>
type Oauth2Token struct {
	AccessToken string    `gorm:"column:access_token;primaryKey" json:"access_token"`
	ClientID    string    `gorm:"column:client_id;not null" json:"client_id"`
	UserID      string    `gorm:"column:user_id;not null" json:"user_id"`
	Scope       string    `gorm:"column:scope;not null" json:"scope"`
	ExpiresAt   time.Time `gorm:"column:expires_at" json:"expires_at"`
	RevokedAt   time.Time `gorm:"column:revoked_at" json:"revoked_at"`
	CreatedAt   time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// TableName Oauth2Token's table name
func (*Oauth2Token) TableName() string {
	return TableNameOauth2Token
}