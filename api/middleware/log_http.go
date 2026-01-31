package middleware

import (
	"log"
	"net/http"
)

type StatusCodeAware struct {
	http.ResponseWriter
	Status int
}

func (s *StatusCodeAware) WriteHeader(statusCode int) {
	s.ResponseWriter.WriteHeader(statusCode)
	s.Status = statusCode
}

func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := &StatusCodeAware{
			ResponseWriter: w,
			Status:         200,
		}
		next.ServeHTTP(response, r)
		log.Printf("[%s] %s %d\n", r.Method, r.URL.Path, response.Status)
	})
}
