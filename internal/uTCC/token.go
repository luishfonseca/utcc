package uTCC

import (
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

func (t Token) ID() int64 {
	return t.id
}

func NewToken(id, n int64, branching int) Token {
	return Token{id, n, time.Now(), branching, n / int64(branching)}
}

func ParseToken(token string, branching int) Token {
	parts := strings.Split(token, "|")
	id, _ := strconv.ParseInt(parts[0], 10, 64)
	n, _ := strconv.ParseInt(parts[1], 10, 64)
	ts, _ := time.Parse(time.RFC3339, parts[2])

	return Token{id, n, ts, branching, n / int64(branching)}
}

func (t *Token) Fraction() Token {
	t.n -= t.fraction
	return Token{t.id, t.fraction, t.ts, t.branching, t.fraction / int64(t.branching)}
}

func (t *Token) Join(other Token) {
	t.n -= other.n
}

func (t Token) Complete() bool {
	return t.n == 0
}

func (t Token) String() string {
	return fmt.Sprintf("%d|%d|%s", t.id, t.n, t.ts.Format(time.RFC3339))
}
