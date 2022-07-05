package handlers

import "net/http"

//ResponseHeader is a middleware handler that adds a header to the response
type ResponseHeader struct {
	handler     http.Handler
	headerName  string
	headerValue string
}

func NewResponseHeaderHandler(handler http.Handler, headerKey, headerValue string) *ResponseHeader {
	return &ResponseHeader{
		handler:     handler,
		headerName:  headerKey,
		headerValue: headerValue,
	}
}

//ServeHTTP handles the request by adding the response header
func (app *ResponseHeader) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//add the header
	w.Header().Add(app.headerName, app.headerValue)
	//call the wrapped handler
	app.handler.ServeHTTP(w, r)
}
