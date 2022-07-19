package hold

type HoldRepository interface {
	SaveHold(hold *Hold) (*Hold, error)
	FindAllHolds() ([]Hold, error)
	FindHoldsFromUser(user string) ([]Hold, error)
}
