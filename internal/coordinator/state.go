package coordinator

import (
	"sync"

	"github.com/luishfonseca/uTCC/internal/uTCC"
)

type Token struct {
	t      uTCC.Token
	tMutex sync.Mutex

	parent     int64
	childCount int
}

type State struct {
	initial   int64
	branching int

	tokens      map[int64]*Token
	tokensMutex sync.RWMutex
}

func int64Pow(base, exp int64) int64 {
	result := int64(1)
	for i := int64(0); i < exp; i++ {
		result *= base
	}
	return result
}

// NewState creates a new state
func NewState(depth, branching int) *State {
	return &State{
		initial:   int64Pow(int64(branching), int64(depth)),
		branching: branching,

		tokens:      make(map[int64]*Token),
		tokensMutex: sync.RWMutex{},
	}
}
