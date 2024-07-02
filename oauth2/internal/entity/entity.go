package entity

type AuthTokens struct {
	AccessToken  string
	RefreshToken string
	Expiry       int64
}

type SessionSigninForm struct {
	Email string
	Error string
}
