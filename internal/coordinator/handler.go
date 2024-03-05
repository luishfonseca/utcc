package coordinator

import (
	"log"
	"math/rand"
	"sync"

	"github.com/valyala/fasthttp"

	"github.com/luishfonseca/uTCC/internal/uTCC"
)

func requestHandler(ctx *fasthttp.RequestCtx, state *State) {
	i := rand.Int63()
	token := Token{
		t:      uTCC.NewToken(i, state.initial, state.branching),
		tMutex: sync.Mutex{},

		childCount: 0,
		parent:     i,
	}

	if len(ctx.Request.Header.Peek("tcc-token")) > 0 {
		parentToken := uTCC.ParseToken(string(ctx.Request.Header.Peek("tcc-token")), state.branching)
		token.parent = parentToken.ID()

		state.tokensMutex.RLock()
		parent, ok := state.tokens[parentToken.ID()]
		state.tokensMutex.RUnlock()

		if !ok {
			log.Fatalf("Parent Token <%d> not found", parentToken.ID())
			ctx.Response.SetStatusCode(fasthttp.StatusBadRequest)
			return
		}

		parent.childCount++

		defer log.Printf("Token <%d> created with parent <%d>", i, parentToken.ID())
	} else {
		defer log.Printf("Root Token <%d> created", i)
	}

	state.tokensMutex.Lock()
	for {
		if _, ok := state.tokens[i]; !ok {
			state.tokens[i] = &token
			break
		}
		i = rand.Int63()

		token.t = uTCC.NewToken(i, state.initial, state.branching)
	}
	state.tokensMutex.Unlock()

	ctx.Response.Header.Set("tcc-token", token.t.String())
	ctx.Response.SetStatusCode(fasthttp.StatusOK)
}

func removeTokenIfComplete(token *Token, state *State) (bool, *Token) {
	if token.t.Complete() && token.childCount == 0 {
		delete(state.tokens, token.t.ID())

		if token.t.ID() == token.parent { // Root token
			return true, token
		}

		if parent, ok := state.tokens[token.parent]; ok {
			parent.childCount--
			return removeTokenIfComplete(parent, state)
		}
	}

	return false, token
}

func returnHandler(ctx *fasthttp.RequestCtx, state *State) {
	returned_token := uTCC.ParseToken(string(ctx.Request.Header.Peek("tcc-token")), state.branching)

	state.tokensMutex.RLock()
	token, ok := state.tokens[returned_token.ID()]
	state.tokensMutex.RUnlock()

	if !ok {
		log.Fatalf("Token <%d> not found", returned_token.ID())
		ctx.Response.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	token.tMutex.Lock()
	token.t.Join(returned_token)
	token.tMutex.Unlock()

	var to_commit *Token
	to_commit = nil
	if token.t.Complete() {
		state.tokensMutex.Lock()
		if root, token := removeTokenIfComplete(token, state); root {
			to_commit = token
		}
		state.tokensMutex.Unlock()
	}

	if to_commit != nil {
		defer log.Printf("COMMIT <%d>", to_commit.t.ID())
	}

	ctx.Response.SetStatusCode(fasthttp.StatusOK)
}

// Handler for the coordinator
func Handler(ctx *fasthttp.RequestCtx, state *State) {
	if string(ctx.Path()) == "/request_token" {
		requestHandler(ctx, state)
		return
	}

	if string(ctx.Path()) == "/return_token" {
		returnHandler(ctx, state)
		return
	}

	log.Fatalf("Unknown call: %s", ctx.Path())
}
