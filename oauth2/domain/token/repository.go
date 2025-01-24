package token

type IAuthorizationRepository interface {
	StoreToken(*Token) error
	FindToken(accessToken string) (*Token, error)
	RevokeToken(accessToken string) error
}
