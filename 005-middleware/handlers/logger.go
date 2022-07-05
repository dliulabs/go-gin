package handlers

import (
	"log"
	"net/http"
	"time"
)

//Logger is a middleware handler that does request logging
type LoggerHandler struct {
	handler http.Handler
}

//NewLogger constructs a new Logger middleware handler
func NewLoggerHandler(h http.Handler) *LoggerHandler {
	return &LoggerHandler{
		handler: h,
	}
}

//ServeHTTP handles the request by passing it to the real
//handler and logging the request details
func (l *LoggerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	l.handler.ServeHTTP(w, r)
	log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
}
