// Handling complex paths with wildcards
// foo/* will match foo/bar but not foo/bar/baz.
// To match foo/bar/baz, we should use foo/*/* (path lib restriction)
package main

import (
	"fmt"
	"net/http"
	"path"
	"strings"
)

func main() {
	pr := newPathResolver()
	pr.Add("GET /hello", hello)
	pr.Add("* /bye/*", bye)
	http.ListenAndServe(":8080", pr)
}

func newPathResolver() *pathResolver {
	return &pathResolver{make(map[string]http.HandlerFunc)}
}

type pathResolver struct {
	handlers map[string]http.HandlerFunc
}

func (p *pathResolver) Add(path string, handler http.HandlerFunc) {
	p.handlers[path] = handler
}

func (p *pathResolver) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	check := req.Method + " " + req.URL.Path
	for pattern, handlerFunc := range p.handlers {
		if ok, err := path.Match(pattern, check); ok && err == nil {
			handlerFunc(res, req)
			return
		} else if err != nil {
			fmt.Fprint(res, err)
		}
	}
	http.NotFound(res, req)
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
