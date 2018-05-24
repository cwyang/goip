// Matching paths to content
package main

import (
	"fmt"
	"net/http"
	"strings"
)

func main() {
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/bye/", bye)
	http.HandleFunc("/", root)
	http.ListenAndServe(":8080", nil)
}

func hello(res http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()
	name := query.Get("name")
	if name == "" {
		name = "cwyang"
	}
	fmt.Fprint(res, "Hi, ", name)
}

func bye(res http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	ps := strings.Split(path, "/")
	name := ps[2]
	if name == "" {		
		name = "cwyang"
	}
	fmt.Fprint(res, "Bye, ", name)
}

func root(res http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		http.NotFound(res, req)
		return
	}
	fmt.Fprint(res, "Ho!")
}
