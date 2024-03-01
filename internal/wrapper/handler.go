package wrapper

import (
	"log"
	"strconv"
	"strings"

	"github.com/valyala/fasthttp"
)

func forward(ctx *fasthttp.RequestCtx, uTCC *State, addr string) {
	ctx.Request.SetHost(addr)
	if err := uTCC.Client().Do(&ctx.Request, &ctx.Response); err != nil {
		log.Fatalf("error in fasthttp.Do: %v", err)
	}
}

func coordHandler(ctx *fasthttp.RequestCtx, uTCC *State) {
	// Handle a call from the coordinator
}

func stateHandler(ctx *fasthttp.RequestCtx, uTCC *State) {
	// Intercepting a state access
}

func invokeHandler(ctx *fasthttp.RequestCtx, uTCC *State) {
	// Parse header to int
	id, _ := strconv.Atoi(string(ctx.Request.Header.Peek("tcc-id")))

	// Get a fraction of the token
	token, err := uTCC.GetTokenFraction(id)
	if err != nil {
		log.Fatalf("error in GetTokenFraction: %v", err)
	}

	ctx.Request.Header.Del("tcc-id")
	ctx.Request.Header.Set("tcc-token", token)

	forward(ctx, uTCC, uTCC.DaprAddr())
}

func daprToAppHandler(ctx *fasthttp.RequestCtx, uTCC *State) {
	id := uTCC.StoreToken(string(ctx.Request.Header.Peek("tcc-token")))

	ctx.Request.Header.Del("tcc-token")
	ctx.Request.Header.Set("tcc-id", strconv.Itoa(id))

	forward(ctx, uTCC, uTCC.AppAddr())

	if !uTCC.HasRemainingToken(id) {
		return
	}

	req := fasthttp.AcquireRequest()
	req.SetRequestURI("http://" + uTCC.CoordAddr() + "/return_token")
	req.Header.SetMethod("POST")
	req.Header.Set("tcc-token", uTCC.GetRemainingToken(id))

	if err := uTCC.Client().Do(req, nil); err != nil {
		log.Fatalf("error in fasthttp.Do: %v", err)
	}

	fasthttp.ReleaseRequest(req)
}

func appToDaprHandler(ctx *fasthttp.RequestCtx, uTCC *State) {
	// Intercepting an invocation
	if strings.HasPrefix(string(ctx.Request.URI().Path()), "/v1.0/invoke/") {
		invokeHandler(ctx, uTCC)
		return
	}

	// Intercepting a state access
	if strings.HasPrefix(string(ctx.Request.URI().Path()), "/v1.0/state/") ||
		strings.HasPrefix(string(ctx.Request.URI().Path()), "/v1.0-alpha/state/") {
		stateHandler(ctx, uTCC)
		return
	}

	log.Fatalf("Unknown App to Dapr call: %s", ctx.Request.URI().Path())
}

func Handler(ctx *fasthttp.RequestCtx, uTCC *State) {
	// Intercepting a call from dapr to the application
	if len(ctx.Request.Header.Peek("tcc-token")) > 0 {
		daprToAppHandler(ctx, uTCC)
		return
	}

	// Intercepting a call from the application to dapr
	if len(ctx.Request.Header.Peek("tcc-id")) > 0 {
		appToDaprHandler(ctx, uTCC)
		return
	}

	// Handle a call from the coordinator
	if strings.HasPrefix(string(ctx.Request.URI().Path()), "/__tcc") {
		coordHandler(ctx, uTCC)
		return
	}

	// Forward dapr internal calls to the application
	if strings.HasPrefix(string(ctx.Request.URI().Path()), "/dapr") {
		forward(ctx, uTCC, uTCC.AppAddr())
		return
	}

	log.Fatalf("Unknown call: %s", ctx.Request.URI().Path())
}
