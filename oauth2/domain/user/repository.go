package user

type UserRepository interface {
	FindUserByEmail(email string) (*User, error)
}
