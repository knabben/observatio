package web

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
)

func WithKubernetes(client, config any) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), "client", client)
			ctx = context.WithValue(ctx, "config", config)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
