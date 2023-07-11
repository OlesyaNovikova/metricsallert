package middleware

import (
	"context"
	"net/http"
)

func WithCtx(ctx context.Context, h http.HandlerFunc) http.HandlerFunc {
	ctxFn := func(w http.ResponseWriter, r *http.Request) {

		r = r.WithContext(ctx)
		// передаём управление хендлеру
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(ctxFn)
}
