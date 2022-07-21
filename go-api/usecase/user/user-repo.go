package user

type UserRepository interface {
	SaveUser(user *User) (*User, error)
	FindAllUsers() ([]User, error)
	FindUserByAddress(address string) (*User, error)
}

type UserERC20Adapter interface {
	GetUserBalance(address string) (*int, error)
	MintTokensToUser(address string, amount int64) error
}
