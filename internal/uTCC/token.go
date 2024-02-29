package uTCC

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Token struct {
	n      int64
	ts     time.Time
	client string
}

func Parse(token string) Token {
	parts := strings.Split(token, ":")
	n, _ := strconv.ParseInt(parts[0], 10, 64)
	ts, _ := time.Parse(time.RFC3339, parts[1])
	client := parts[2]

	return Token{n, ts, client}
}

func (t *Token) Fraction(branching int) Token {
	fraction := t.n / int64(branching)
	t.n -= fraction
	return Token{fraction, t.ts, t.client}
}

func (t Token) String() string {
	return fmt.Sprintf("%d:%s:%s", t.n, t.ts.Format(time.RFC3339), t.client)
}
