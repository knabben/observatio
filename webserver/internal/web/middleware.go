package web

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func WithKubernetes(client, config any) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), "client", client) // nolint
			ctx = context.WithValue(ctx, "config", config)          // nolint
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func WithLogger() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.FromContext(r.Context()).WithValues(
				"url", r.URL.String(),
				"method", r.Method,
				"remote_addr", r.RemoteAddr,
			).Info("Request")
			next.ServeHTTP(w, r)
		})
	}
}
