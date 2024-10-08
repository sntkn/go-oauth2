// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameOauth2RefreshToken = "oauth2_refresh_tokens"

// Oauth2RefreshToken mapped from table <oauth2_refresh_tokens>
type Oauth2RefreshToken struct {
	RefreshToken string    `gorm:"column:refresh_token;primaryKey" json:"refresh_token"`
	AccessToken  string    `gorm:"column:access_token;not null" json:"access_token"`
	ExpiresAt    time.Time `gorm:"column:expires_at" json:"expires_at"`
	CreatedAt    time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// TableName Oauth2RefreshToken's table name
func (*Oauth2RefreshToken) TableName() string {
	return TableNameOauth2RefreshToken
}
