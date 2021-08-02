package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/JIexa24/chef-webapi/httpserver/errors"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

// IDContext set the id context value.
func IDContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		id := chi.URLParam(request, "id")
		if id == "" {
			_ = render.Render(writer, request,
				errors.ErrorInvalidRequest(fmt.Errorf("id must be set")))
			return
		}
		ctx := context.WithValue(request.Context(), idKey, id)
		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}
