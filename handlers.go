package main

import(
	"fmt"
	"encoding/json"

	"github.com/valyala/fasthttp"
)

func Index(ctx *fasthttp.RequestCtx) {
	fmt.Fprintln(ctx, "Welcom!")
}

func Hello(ctx *fasthttp.RequestCtx) {
	r := doRequest(ctx)
	encode := json.NewEncoder(ctx)
	encode.Encode(&r)
	ctx.SetStatusCode(fasthttp.StatusOK)
}