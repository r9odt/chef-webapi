package middleware

import (
	"fmt"
	"net/http"

	web "github.com/JIexa24/chef-webapi"

	"github.com/JIexa24/chef-webapi/httpserver/errors"

	"github.com/go-chi/render"
)

// IsAuth check for user has session for access to api.
func IsAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		session := request.Header.Get(SessionHeader)
		auth := request.Header.Get(AuthHeader)
		responseError := fmt.Errorf("access without session")
		if auth == "X-Auth-Header" {
		} else if session == "" || session == "null" {
			web.App.Logger.Error(responseError)
			_ = render.Render(writer, request,
				errors.ErrorUnauthorized(responseError))
			return
		}
		s, err := web.App.DB.GetSessionByUUID(session)
		if auth == "X-Auth-Header" {
		} else if s == nil || err != nil {
			web.App.Logger.Error(responseError)
			_ = render.Render(writer, request,
				errors.ErrorUnauthorized(responseError))
			return
		}
		next.ServeHTTP(writer, request)
	})
}
