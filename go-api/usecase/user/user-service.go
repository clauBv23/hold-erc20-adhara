package user

import (
	"errors"
	"math/rand"
)

type User struct {
	Id      int64
	Address string
}

type UserService interface {
	CreateUser(hold *User) (*User, error)
	FindAllUsers() ([]User, error)
	FindUserBalance(address string) (*int, error)
	FindUserByAddress(address string) (*User, error)
}

type service struct {
	repo    UserRepository
	adapter UserERC20Adapter
}

func NewUserService(repo UserRepository, adapter UserERC20Adapter) UserService {
	return &service{repo: repo, adapter: adapter}
}

func (s *service) CreateUser(hold *User) (*User, error) {

	user, err := s.FindUserByAddress(hold.Address)
	if err != nil {
		return nil, err
	}
	if user.Id != 0 {
		err := errors.New("The user is already registered")
		return nil, err
	}

	hold.Id = rand.Int63()

	// mint tokens to the registered account
	err = s.adapter.MintTokensToUser(hold.Address, 500)
	if err != nil {
		return nil, err
	}

	return s.repo.SaveUser(hold)
}

func (s *service) FindAllUsers() ([]User, error) {
	return s.repo.FindAllUsers()
}

func (s *service) FindUserByAddress(address string) (*User, error) {
	return s.repo.FindUserByAddress(address)
}

func (s *service) FindUserBalance(address string) (*int, error) {
	// todo avoid calling the smart contract if the user is not registered
	// todo change the request to make them with id

	return s.adapter.GetUserBalance(address)
}
