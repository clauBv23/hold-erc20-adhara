package hold

import (
	"cleanGo/api/usecase/user"
	"errors"
	"log"
	"math/rand"
)

type Hold struct {
	Id     int64
	Amount int64
	User   string
}

type HoldService interface {
	ValidateHold(hold *Hold) error
	CreateHold(hold *Hold) (*Hold, error)
	FindAllHolds() ([]Hold, error)
	FindHoldsFromUser(user string) ([]Hold, error)
}

type service struct {
	repo        HoldRepository
	userService user.UserService
}

func NewHoldService(repo HoldRepository, userService user.UserService) HoldService {
	return &service{repo: repo, userService: userService}
}

func (s *service) ValidateHold(hold *Hold) error {

	if hold == nil {
		err := errors.New("The hold is empty.")
		return err
	}

	if hold.Amount <= 0 {
		err := errors.New("The hold amount must be bigger than zero.")
		return err
	}

	user, err := s.userService.FindUserByAddress(hold.User)
	if err != nil {
		log.Fatalf("Failed to find user: %v", err)
		return err
	}

	if user.Id == 0 {
		err := errors.New("The hold user is not registered.")
		return err
	}
	return nil
}

func (s *service) CreateHold(hold *Hold) (*Hold, error) {
	//todo create the hold on the smart contract
	//todo add the flow when there are more than 5 registered holds to select a random winner
	//			todo a status field on the hold may be needed
	err := s.ValidateHold(hold)
	if err != nil {
		return nil, err
	}
	hold.Id = rand.Int63()
	return s.repo.SaveHold(hold)
}

func (s *service) FindAllHolds() ([]Hold, error) {
	return s.repo.FindAllHolds()
}

func (s *service) FindHoldsFromUser(user string) ([]Hold, error) {
	return s.repo.FindHoldsFromUser(user)

}
