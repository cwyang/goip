// Configuration via environment variables
package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/shutdown", shutdown) // Do NOT do in production code
	http.HandleFunc("/", homePage)
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}

func shutdown(res http/ResponseWriter, req *http.Request) {	// antipattern
	os.Exit(0)
}

func homePage(res http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		http.NotFound(res, req)
		return
	}
	fmt.Fprint(res, "Hello, world!")
}
