package chariot

import (
	"net/http"
)

func HTTPHandler(r func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		r(res, req)
	})
}
