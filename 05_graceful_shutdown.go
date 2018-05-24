// Graceful shutdowns using manners
package main

import (
	"github.com/braintree/manners" // same interface for ListenAndServe

	"fmt"
	"net/http"
	"os"
	"os/signal"
)

func main() {
	handler := newHandler()

	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, os.Kill)
	go listenForShutDown(ch)

	manners.ListenAndServe(":8080", handler)
}

func newHandler() *handler {
	return &handler{}
}

type handler struct{}

func (h *handler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()
	name := query.Get("name")
	if name == "" {
		name = "cwyang"
	}
	fmt.Fprint(res, "Hi, ", name)
}

func listenForShutDown(ch <-chan os.Signal) {
	<-ch
	manners.Close()
}
