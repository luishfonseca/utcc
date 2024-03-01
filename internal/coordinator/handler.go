package coordinator

import (
	"log"

	"github.com/valyala/fasthttp"
)

func requestHandler(ctx *fasthttp.RequestCtx, state *State) {
	// Request a token
}

func returnHandler(ctx *fasthttp.RequestCtx, state *State) {
	// Give back the token
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
