package token

type TokenRepository interface {
	StoreToken(*Token) error
	FindToken(accessToken string) (*Token, error)
	RevokeToken(accessToken string) error
}
