package authorization_code

import (
	"time"
)

type AuthorizationCodeRepository interface {
	FindAuthorizationCode(string) (*AuthorizationCode, error)
	StoreAuthorizationCode(*AuthorizationCode) error
	FindValidAuthorizationCode(string, time.Time) (*AuthorizationCode, error)
	RevokeCode(code string) error
}
