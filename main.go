package main

import(
	"log"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

func main() {
	router := fasthttprouter.New()
	router.GET("/", Index)
	router.GET("/par/est", Hello)

	log.Fatal(fasthttp.ListenAndServe(":8080", router.Handler))
}