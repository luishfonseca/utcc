package wrapper

import (
	"log"
	"strconv"
	"strings"

	"github.com/valyala/fasthttp"
)

func forward(ctx *fasthttp.RequestCtx, state *State, addr string) {
	ctx.Request.SetHost(addr)
	if err := state.Client().Do(&ctx.Request, &ctx.Response); err != nil {
		log.Fatalf("error in fasthttp.Do: %v", err)
	}
}

func coordHandler(ctx *fasthttp.RequestCtx, state *State) {
	// Handle a call from the coordinator
}

func stateHandler(ctx *fasthttp.RequestCtx, state *State) {
	// Intercepting a state access
}

func invokeHandler(ctx *fasthttp.RequestCtx, state *State) {
	// Parse header to int
	id, _ := strconv.Atoi(string(ctx.Request.Header.Peek("tcc-id")))

	// Get a fraction of the token
	token, err := state.GetTokenFraction(id, func(prev string) string {
		req := fasthttp.AcquireRequest()
		req.SetRequestURI("http://" + state.CoordAddr() + "/request_token")
		req.Header.SetMethod("POST")
		req.Header.Set("tcc-token", prev)

		resp := fasthttp.AcquireResponse()
		if err := state.Client().Do(req, resp); err != nil {
			log.Fatalf("error in fasthttp.Do: %v", err)
		}

		fasthttp.ReleaseRequest(req)

		return string(resp.Header.Peek("tcc-token"))
	})

	if err != nil {
		log.Fatalf("error in GetTokenFraction: %v", err)
	}

	ctx.Request.Header.Del("tcc-id")
	ctx.Request.Header.Set("tcc-token", token)

	forward(ctx, state, state.DaprAddr())
}

func daprToAppHandler(ctx *fasthttp.RequestCtx, state *State) {
	id := state.StoreToken(string(ctx.Request.Header.Peek("tcc-token")))

	ctx.Request.Header.Del("tcc-token")
	ctx.Request.Header.Set("tcc-id", strconv.Itoa(id))

	forward(ctx, state, state.AppAddr())

	go func() {
		if !state.HasRemainingToken(id) {
			return
		}

		req := fasthttp.AcquireRequest()
		req.SetRequestURI("http://" + state.CoordAddr() + "/return_token")
		req.Header.SetMethod("POST")
		req.Header.Set("tcc-token", state.GetRemainingToken(id))

		if err := state.Client().Do(req, nil); err != nil {
			log.Fatalf("error in fasthttp.Do: %v", err)
		}

		fasthttp.ReleaseRequest(req)
	}()
}

func appToDaprHandler(ctx *fasthttp.RequestCtx, state *State) {
	// Intercepting an invocation
	if strings.HasPrefix(string(ctx.Request.URI().Path()), "/v1.0/invoke/") {
		invokeHandler(ctx, state)
		return
	}

	// Intercepting a state access
	if strings.HasPrefix(string(ctx.Request.URI().Path()), "/v1.0/state/") ||
		strings.HasPrefix(string(ctx.Request.URI().Path()), "/v1.0-alpha/state/") {
		stateHandler(ctx, state)
		return
	}

	log.Fatalf("Unknown App to Dapr call: %s", ctx.Request.URI().Path())
}

func Handler(ctx *fasthttp.RequestCtx, state *State) {
	// Intercepting a call from dapr to the application
	if len(ctx.Request.Header.Peek("tcc-token")) > 0 {
		daprToAppHandler(ctx, state)
		return
	}

	// Intercepting a call from the application to dapr
	if len(ctx.Request.Header.Peek("tcc-id")) > 0 {
		appToDaprHandler(ctx, state)
		return
	}

	// Handle a call from the coordinator
	if strings.HasPrefix(string(ctx.Request.URI().Path()), "/__tcc") {
		coordHandler(ctx, state)
		return
	}

	// Forward dapr internal calls to the application
	if strings.HasPrefix(string(ctx.Request.URI().Path()), "/dapr") {
		forward(ctx, state, state.AppAddr())
		return
	}

	log.Fatalf("Unknown call: %s", ctx.Request.URI().Path())
}
