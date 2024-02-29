package uTCC

import (
	"math/rand"

	"github.com/valyala/fasthttp"
)

// State of the uTCC
type State struct {
	client *fasthttp.Client

	dapr_addr  string
	app_addr   string
	coord_addr string

	branching     int
	active_tokens map[int]Token
}

func (s *State) Client() *fasthttp.Client {
	return s.client
}

func (s State) DaprAddr() string {
	return s.dapr_addr
}

func (s State) AppAddr() string {
	return s.app_addr
}

func (s State) CoordAddr() string {
	return s.coord_addr
}

// NewState creates a new state
func NewState(dapr_addr, app_addr, coord_addr string, branching int) *State {
	return &State{
		client:        &fasthttp.Client{},
		dapr_addr:     dapr_addr,
		app_addr:      app_addr,
		coord_addr:    coord_addr,
		branching:     branching,
		active_tokens: make(map[int]Token),
	}
}

// Store a token in the state
func (s *State) StoreToken(token string) int {
	i := rand.Int()
	for {
		if _, ok := s.active_tokens[i]; !ok {
			break
		}
		i = rand.Int()
	}

	s.active_tokens[i] = Parse(token)

	return i
}

// Get a token fraction from the state
func (s *State) GetTokenFraction(id int) string {
	token := s.active_tokens[id]
	fraction := token.Fraction(s.branching)
	s.active_tokens[id] = token
	return fraction.String()
}

// Get remaining token from the state
func (s *State) GetRemainingToken(id int) string {
	token := s.active_tokens[id]
	delete(s.active_tokens, id)
	return token.String()
}
