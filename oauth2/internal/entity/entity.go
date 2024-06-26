package entity

type AuthTokens struct {
	AccessToken  string
	RefreshToken string
	Expiry       int64
}
