package refresh_token

import (
	"time"
)

type RefreshToken struct {
	RefreshToken string
	AccessToken  string
	ExpiresAt    time.Time
}
