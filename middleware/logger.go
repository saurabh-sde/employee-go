package middleware

import (
	"net/http"
	"time"

	"github.com/saurabh-sde/employee-go/utility"
)

// Middleware for logging API methods and path
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		utility.Print("API: BEGIN ", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		utility.Print("API: DONE ", r.Method, r.URL.Path, time.Since(time.Now()).String())
	})
}
