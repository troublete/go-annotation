package home

import (
	"fmt"
	"net/http"
)

// Home is the default route, triggered when root is called
//
// chariot.route{path=/{name}}
func Home(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, "hello %v!", r.PathValue("name"))
}

// Another is another route defined
//
// chariot.route{path=/special/{param}/name}
func Another(w http.ResponseWriter, _ *http.Request) {
	_, _ = fmt.Fprintf(w, "hello world")
}
