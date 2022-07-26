package hold

import (
	"cleanGo/api/usecase/user"
	"errors"
	"log"
)

type Hold struct {
	Id     int64
	Amount int64
	User   string
	Status string
}

type HoldService interface {
	ValidateHold(hold *Hold) error
	CreateHold(hold *Hold) (*Hold, error)
	FindAllHolds() ([]Hold, error)
	FindHoldsFromUser(user string) ([]Hold, error)
}

type service struct {
	adapter     HoldERC20Adapter
	repo        HoldRepository
	userService user.UserService
}

func NewHoldService(repo HoldRepository, userService user.UserService, adapter HoldERC20Adapter) HoldService {
	return &service{repo: repo, userService: userService, adapter: adapter}
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
	//todo add the flow when there are more than 5 registered holds to select a random winner
	//			todo a status field on the hold may be needed
	err := s.ValidateHold(hold)
	if err != nil {
		return nil, err
	}

	// create the hold on the smart contract
	id, err := s.adapter.CreateHold(hold.User, hold.Amount)
	if err != nil {
		return nil, err
	}
	hold.Id = int64(id)
	hold.Status = "CREATED"

	newHold, err := s.repo.SaveHold(hold)
	// todo  check if there are 5 created and execute them

	return s.checkIfHaveExecHolds(newHold)
}

func (s *service) checkIfHaveExecHolds(newHold *Hold) (*Hold, error) {
	//	get the holds on create
	// if there are 5 execute the holds
	
	return newHold, nil
}

func (s *service) FindAllHolds() ([]Hold, error) {
	return s.repo.FindAllHolds()
}

func (s *service) FindAllHoldsOnCreated() ([]Hold, error) {
	return s.repo.FindAllHoldsOnCreated()
}

func (s *service) FindHoldsFromUser(user string) ([]Hold, error) {
	return s.repo.FindHoldsFromUser(user)

}
