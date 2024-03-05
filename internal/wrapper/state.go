package wrapper

import (
	"errors"
	"log"
	"math/rand"

	"github.com/valyala/fasthttp"

	"github.com/luishfonseca/utcc/internal/uTCC"
)

// State of the uTCC
type State struct {
	client *fasthttp.Client

	dapr_addr  string
	app_addr   string
	coord_addr string

	branching     int
	active_tokens map[int]uTCC.Token
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
		active_tokens: make(map[int]uTCC.Token),
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

	s.active_tokens[i] = uTCC.ParseToken(token, s.branching)

	log.Printf("Stored token <%d>: %d", i, s.active_tokens[i].N())

	return i
}

// Get a token fraction from the state
func (s *State) GetTokenFraction(id int, request func(string) string) (string, error) {
	token, ok := s.active_tokens[id]
	if !ok {
		return "", errors.New("Token not found")
	}

	if token.Complete() {
		token = uTCC.ParseToken(request(token.String()), s.branching)
		log.Printf("Token <%d> is fully consumed, requested new token: <%d>", id, token.N())
	}

	fraction := token.Fraction()
	s.active_tokens[id] = token

	log.Printf("Fraction of token <%d>: %d (remaining: %d)", id, fraction.N(), token.N())

	return fraction.String(), nil
}

func (s *State) HasRemainingToken(id int) bool {
	token, ok := s.active_tokens[id]
	return ok && !token.Complete()
}

// Get remaining token from the state
func (s *State) GetRemainingToken(id int) string {
	token := s.active_tokens[id]
	log.Printf("Remaining token <%d>: %d", id, token.N())

	delete(s.active_tokens, id)
	return token.String()
}
