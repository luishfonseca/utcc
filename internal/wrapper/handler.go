package wrapper

import (
	"log"
	"strconv"
	"strings"

	"github.com/valyala/fasthttp"

	"github.com/luishfonseca/uTCC/internal/uTCC"
)

func invokeHandler(ctx *fasthttp.RequestCtx, uTCC *uTCC.State) {
	// Parse header to int
	id, _ := strconv.Atoi(string(ctx.Request.Header.Peek("tcc-id")))

	// Get a fraction of the token
	token := uTCC.GetTokenFraction(id)

	// Build the request to Dapr
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(uTCC.DaprAddr() + string(ctx.Request.URI().Path()))
	req.Header.SetMethodBytes(ctx.Request.Header.Method())
	req.Header.SetContentTypeBytes(ctx.Request.Header.ContentType())
	req.Header.Set("tcc-token", token)
	req.SetBody(ctx.Request.Body())

	// Send the request to Dapr
	resp := fasthttp.AcquireResponse()
	if err := uTCC.Client().Do(req, resp); err != nil {
		log.Fatalf("error in fasthttp.Do: %v", err)
	}

	// Build the response to the application
	ctx.Response.SetStatusCode(resp.StatusCode())
	ctx.Response.Header.SetContentTypeBytes(resp.Header.ContentType())
	ctx.Response.SetBody(resp.Body())

	// Release the request and response
	fasthttp.ReleaseRequest(req)
	fasthttp.ReleaseResponse(resp)
}

func stateHandler(ctx *fasthttp.RequestCtx, uTCC *uTCC.State) {
	// Intercepting a state access
}

func daprToAppHandler(ctx *fasthttp.RequestCtx, uTCC *uTCC.State) {
	id := uTCC.StoreToken(string(ctx.Request.Header.Peek("tcc-token")))

	// Build the request to the application
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(uTCC.AppAddr() + string(ctx.Request.URI().Path()))
	req.Header.SetMethodBytes(ctx.Request.Header.Method())
	req.Header.SetContentTypeBytes(ctx.Request.Header.ContentType())
	req.Header.Set("tcc-id", strconv.Itoa(id))
	req.SetBody(ctx.Request.Body())

	// Send the request to the application
	resp := fasthttp.AcquireResponse()
	if err := uTCC.Client().Do(req, resp); err != nil {
		log.Fatalf("error in fasthttp.Do: %v", err)
	}

	// Build the response to Dapr
	ctx.Response.SetStatusCode(resp.StatusCode())
	ctx.Response.Header.SetContentTypeBytes(resp.Header.ContentType())
	ctx.Response.SetBody(resp.Body())

	// Build the request to the coordinator
	reqCoord := fasthttp.AcquireRequest()
	reqCoord.SetRequestURI(uTCC.CoordAddr() + "/give_back")
	reqCoord.Header.SetMethod("POST")
	reqCoord.Header.Set("tcc-token", uTCC.GetRemainingToken(id))

	// Send the request to the coordinator
	if err := uTCC.Client().Do(reqCoord, nil); err != nil {
		log.Fatalf("error in fasthttp.Do: %v", err)
	}
}

func coordHandler(ctx *fasthttp.RequestCtx, uTCC *uTCC.State) {
	// Handle a call from the coordinator
}

func appToDaprHandler(ctx *fasthttp.RequestCtx, uTCC *uTCC.State) {
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
}

func Handler(ctx *fasthttp.RequestCtx, uTCC *uTCC.State) {
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
	if strings.HasPrefix(string(ctx.Request.URI().Path()), "/__tcc/") {
		coordHandler(ctx, uTCC)
		return
	}

	log.Printf("Invalid request")
}
