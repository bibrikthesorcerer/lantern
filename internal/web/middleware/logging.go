package middleware

import (
	"net/http"
	"time"

	clog "github.com/charmbracelet/log"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		received := time.Now()
		next.ServeHTTP(w, r)
		clog.Printf("%v %v %v", r.Method, r.URL.Path, time.Since(received))
	})
}
