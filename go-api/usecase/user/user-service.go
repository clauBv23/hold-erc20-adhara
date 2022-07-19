package user

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
)

type User struct {
	Id      int64
	Address string
}

type UserService interface {
	CreateUser(hold *User) (*User, error)
	FindAllUsers() ([]User, error)
	FindUserBalance(userId int64) (int, error)
	FindUserByAddress(address string) (*User, error)
}

type service struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) UserService {
	return &service{repo: repo}
}

func (s *service) CreateUser(hold *User) (*User, error) {
	//todo register the user on the smart contract

	user, err := s.FindUserByAddress(hold.Address)
	if err != nil {
		log.Fatalf("Failed to get the user: %v", err)
		return nil, err
	}
	if user.Id != 0 {
		err := errors.New("The user is already registered")
		return nil, err
	}

	hold.Id = rand.Int63()
	return s.repo.SaveUser(hold)
}

func (s *service) FindAllUsers() ([]User, error) {
	return s.repo.FindAllUsers()
}

func (s *service) FindUserByAddress(address string) (*User, error) {
	return s.repo.FindUserByAddress(address)
}

func (s *service) FindUserBalance(userId int64) (int, error) {
	// todo get the user balance from the smart contract
	// todo avoid calling the smart contract if the user is not registered
	fmt.Println(userId)
	return 0, nil
}
