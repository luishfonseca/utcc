package main

import (
	"flag"
	"log"

	"github.com/valyala/fasthttp"

	"github.com/luishfonseca/uTCC/internal/uTCC"
	"github.com/luishfonseca/uTCC/internal/wrapper"
)

var (
	wrapper_addr = flag.String("wrapper", "localhost:3500", "Wrapper listener address")

	app_addr   = flag.String("app", "localhost:3501", "App address")
	dapr_addr  = flag.String("dapr", "localhost:3502", "Dapr address")
	coord_addr = flag.String("coord", "localhost:3503", "Coordinator address")

	branching = flag.Int("branching", 4, "Branching factor")
)

func main() {
	flag.Parse()
	uTCC := uTCC.NewState(*dapr_addr, *app_addr, *coord_addr, *branching)

	go func() {
		if err := fasthttp.ListenAndServe(*wrapper_addr, func(ctx *fasthttp.RequestCtx) { wrapper.Handler(ctx, uTCC) }); err != nil {
			log.Fatalf("error in ListenAndServe: %v", err)
		}
	}()

	log.Printf("uTCC is running")

	select {}
}