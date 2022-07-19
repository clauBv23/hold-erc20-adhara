package user

type UserRepository interface {
	SaveUser(user *User) (*User, error)
	FindAllUsers() ([]User, error)
	FindUserByAddress(address string) (*User, error)
}
