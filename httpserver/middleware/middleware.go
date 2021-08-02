package middleware

import (
	"net/http"
)

// Middleware package must implement middleware logic and extract data from
// request.

// ContextKey is base key type.
type ContextKey string

func (key ContextKey) String() string {
	return "context key " + string(key)
}

// SessionHeader is header for session key.
var SessionHeader = http.CanonicalHeaderKey("X-Session-Key")

// AuthHeader is header for auth process.
var AuthHeader = http.CanonicalHeaderKey("X-Auth-Header")

var (
	searchIndexKey ContextKey = "searchIndex"
	searchQueryKey ContextKey = "searchQuery"
	idKey          ContextKey = "id"
)

// GetID get ID from request context,
// which was sets in IDContext middleware.
func GetID(request *http.Request) string {
	id := request.Context().Value(idKey).(string)
	return id
}

// GetSearchIndex get search type string from request context,
// which was sets in SearchIndexContext middleware.
func GetSearchIndex(request *http.Request) string {
	return request.Context().Value(searchIndexKey).(string)
}

// GetSearchQuery get search parameter string from request context,
// which was sets in SearchQueryContext middleware.
func GetSearchQuery(request *http.Request) string {
	return request.Context().Value(searchQueryKey).(string)
}
