package middleware

import (
	"log"
	"net/http"
	"time"
)

type Middleware struct{}

func (m *Middleware) LoggerMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		end := time.Now()
		log.Printf("URL: %s, METHOD: %s, Duration: %d\n", r.URL.String(), r.Method, end.Sub(start))
	}
	return http.HandlerFunc(fn)
}

func (m *Middleware) RecoverMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("URL: %s, METHOD: %s, Error: %+v\n", r.URL.String(), r.Method, err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
