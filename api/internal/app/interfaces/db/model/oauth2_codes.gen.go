// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameOauth2Code = "oauth2_codes"

// Oauth2Code mapped from table <oauth2_codes>
type Oauth2Code struct {
	Code        string    `gorm:"column:code;primaryKey" json:"code"`
	ClientID    string    `gorm:"column:client_id;not null" json:"client_id"`
	UserID      string    `gorm:"column:user_id;not null" json:"user_id"`
	Scope       string    `gorm:"column:scope;not null" json:"scope"`
	RedirectURI string    `gorm:"column:redirect_uri;not null" json:"redirect_uri"`
	ExpiresAt   time.Time `gorm:"column:expires_at" json:"expires_at"`
	RevokedAt   time.Time `gorm:"column:revoked_at" json:"revoked_at"`
	CreatedAt   time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// TableName Oauth2Code's table name
func (*Oauth2Code) TableName() string {
	return TableNameOauth2Code
}