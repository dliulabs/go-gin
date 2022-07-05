package main

import (
	"log"
	"net/http"

	"hello/handlers"
)

var helloApp *handlers.HelloHandler
var loggerApp *handlers.LoggerHandler
var headerApp *handlers.ResponseHeader

func init() {
	helloApp = handlers.NewHelloHandler("Hello")

	//wrap entire mux with logger middleware
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/hello", helloApp.HelloHandler)
	mux.HandleFunc("/v1/time", helloApp.CurrentTimeHandler)

	headerApp = handlers.NewResponseHeaderHandler(mux, "X-My-Header", "my header value")
	loggerApp = handlers.NewLoggerHandler(headerApp)
}

func main() {
	addr := "0.0.0.0:9090"
	log.Printf("server is listening at %s", addr)
	//use wrappedMux instead of mux as root handler
	log.Fatal(http.ListenAndServe(addr, loggerApp))
}
