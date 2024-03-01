package main

import (
	"flag"
	"log"

	"github.com/valyala/fasthttp"

	"github.com/luishfonseca/uTCC/internal/coordinator"
)

var (
	coord_addr = flag.String("coord", "localhost:3503", "Coordinator address")

	branching = flag.Int("branching", 4, "Branching factor")
)

func main() {
	flag.Parse()
	coordState := coordinator.NewState(*branching)

	go func() {
		if err := fasthttp.ListenAndServe(*coord_addr, func(ctx *fasthttp.RequestCtx) { coordinator.Handler(ctx, coordState) }); err != nil {
			log.Fatalf("error in ListenAndServe: %v", err)
		}
	}()

	log.Printf("uTCC Coordinator is running on %s", *coord_addr)

	select {}
}
