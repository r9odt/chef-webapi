package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/r9odt/chef-webapi/httpserver/errors"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

// SearchIndexContext set the searchIndex context value.
func SearchIndexContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		searchIndex := chi.URLParam(request, "searchIndex")
		if searchIndex == "" {
			_ = render.Render(writer, request,
				errors.ErrorInvalidRequest(fmt.Errorf("searchIndex must be set")))
			return
		}
		ctx := context.WithValue(request.Context(), searchIndexKey, searchIndex)
		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}

// SearchQueryContext set the searchQuery context value.
func SearchQueryContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		searchQuery := chi.URLParam(request, "searchQuery")
		if searchQuery == "" {
			_ = render.Render(writer, request,
				errors.ErrorInvalidRequest(fmt.Errorf("searchQuery must be set")))
			return
		}
		ctx := context.WithValue(request.Context(), searchQueryKey, searchQuery)
		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}
