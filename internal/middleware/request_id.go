package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type contextKey string

const RequestIDKey contextKey = "request_id"

func RequestID() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := fmt.Sprintf("%d", time.Now().UnixNano())

			ctx := context.WithValue(r.Context(), RequestIDKey, id)
			w.Header().Set("X-Request-ID", id)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
