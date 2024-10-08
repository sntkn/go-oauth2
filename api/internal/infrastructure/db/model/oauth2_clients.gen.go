// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameOauth2Client = "oauth2_clients"

// Oauth2Client mapped from table <oauth2_clients>
type Oauth2Client struct {
	ID           string    `gorm:"column:id;primaryKey;default:gen_random_uuid()" json:"id"`
	Name         string    `gorm:"column:name;not null" json:"name"`
	RedirectUris string    `gorm:"column:redirect_uris;not null" json:"redirect_uris"`
	CreatedAt    time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// TableName Oauth2Client's table name
func (*Oauth2Client) TableName() string {
	return TableNameOauth2Client
}
