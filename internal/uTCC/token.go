package uTCC

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Token struct {
	id int64
	n  int64
	ts time.Time

	branching int
	fraction  int64
}

func (t Token) N() int64 {
	return t.n
}

func ParseToken(token string, branching int) Token {
	parts := strings.Split(token, "|")
	id, _ := strconv.ParseInt(parts[0], 10, 64)
	n, _ := strconv.ParseInt(parts[1], 10, 64)
	ts, _ := time.Parse(time.RFC3339, parts[2])

	return Token{id, n, ts, branching, n / int64(branching)}
}

func (t *Token) Fraction() (Token, error) {
	if t.n == 0 {
		return Token{}, errors.New("Token is already fully branched")
	}

	t.n -= t.fraction
	return Token{t.id, t.fraction, t.ts, t.branching, t.fraction / int64(t.branching)}, nil
}

func (t Token) String() string {
	return fmt.Sprintf("%d|%d|%s", t.n, t.id, t.ts.Format(time.RFC3339))
}
