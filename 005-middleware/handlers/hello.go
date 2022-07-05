package handlers

import (
	"fmt"
	"net/http"
	"time"
)

type HelloHandler struct {
	message string
}

func NewHelloHandler(greeting string) *HelloHandler {
	return &HelloHandler{
		message: greeting,
	}
}

func (app *HelloHandler) HelloHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(app.message))
}

func (app *HelloHandler) CurrentTimeHandler(w http.ResponseWriter, r *http.Request) {
	curTime := time.Now().Format(time.Kitchen)
	w.Write([]byte(fmt.Sprintf("%v: the current time is %v", app.message, curTime)))
}
