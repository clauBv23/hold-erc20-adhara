package hold

type HoldRepository interface {
	SaveHold(hold *Hold) (*Hold, error)
	FindAllHolds() ([]Hold, error)
	FindAllHoldsOnCreated() ([]Hold, error)
	FindHoldsFromUser(user string) ([]Hold, error)
}

type HoldERC20Adapter interface {
	CreateHold(holderAddr string, amount int64) (int64, error)
}
