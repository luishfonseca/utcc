package coordinator

type State struct {
	branching int
}

// NewState creates a new state
func NewState(branching int) *State {
	return &State{
		branching: branching,
	}
}
