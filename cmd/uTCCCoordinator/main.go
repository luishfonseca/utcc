package main

import (
	"flag"
	"log"

	"github.com/valyala/fasthttp"

	"github.com/luishfonseca/utcc/internal/coordinator"
)

var (
	coord_addr = flag.String("coord", "localhost:3503", "Coordinator address")

	branching = flag.Int("branching", 4, "Branching factor")
	depth     = flag.Int("depth", 4, "Depth of the tree")
)

func main() {
	flag.Parse()
	state := coordinator.NewState(*depth, *branching)

	go func() {
		if err := fasthttp.ListenAndServe(*coord_addr, func(ctx *fasthttp.RequestCtx) { coordinator.Handler(ctx, state) }); err != nil {
			log.Fatalf("error in ListenAndServe: %v", err)
		}
	}()

	log.Printf("uTCC Coordinator is running on %s", *coord_addr)

	select {}
}
