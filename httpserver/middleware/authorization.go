package middleware

import (
	"fmt"
	"net/http"

	web "github.com/r9odt/chef-webapi"

	"github.com/r9odt/chef-webapi/httpserver/errors"

	"github.com/go-chi/render"
)

// hasAdminGrants check if user is admin.
func hasAdminGrants(request *http.Request) (bool, int, error) {
	errorCode := 200
	session := request.Header.Get(SessionHeader)
	s, err := web.App.DB.GetSessionByUUID(session)
	if s != nil && err == nil {
		user, err := web.App.DB.GetUserByUsername(s.Username)
		if user != nil && err == nil {
			if !user.IsAdmin {
				return false, errorCode, err
			}
			return true, errorCode, err
		}
		if err != nil {
			errorCode = 500
			return false, errorCode, err
		}
	}
	errorCode = 401
	return false, errorCode, err
}

// HasAccessToUsers check client on aceess to users endpoint.
func HasAccessToUsers(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		responseError := fmt.Errorf("not found")
		isAdmin, code, err := hasAdminGrants(request)
		if code == 401 {
				_ = render.Render(writer, request,
					errors.ErrorUnauthorized(responseError))
				return
		}
		if code == 500 {
				web.App.Logger.Error(err.Error())
				_ = render.Render(writer, request,
					errors.ErrorInternalServer(responseError))
				return
		}
		if !isAdmin {
				web.App.Logger.Error(responseError.Error())
				_ = render.Render(writer, request,
					errors.ErrorNotFound(responseError))
				return
		}
		next.ServeHTTP(writer, request)
	})
}

// HasAccessToModules check client on aceess to users endpoint.
func HasAccessToModules(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		responseError := fmt.Errorf("not found")
		isAdmin, code, err := hasAdminGrants(request)
		if code == 401 {
				_ = render.Render(writer, request,
					errors.ErrorUnauthorized(responseError))
				return
		}
		if code == 500 {
				web.App.Logger.Error(err.Error())
				_ = render.Render(writer, request,
					errors.ErrorInternalServer(responseError))
				return
		}
		if !isAdmin {
				web.App.Logger.Error(responseError.Error())
				_ = render.Render(writer, request,
					errors.ErrorNotFound(responseError))
				return
		}
		next.ServeHTTP(writer, request)
	})
}


// HasAccessToKeys check client on aceess to keys endpoint.
func HasAccessToKeys(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		responseError := fmt.Errorf("not found")
		isAdmin, code, err := hasAdminGrants(request)
		if code == 401 {
				_ = render.Render(writer, request,
					errors.ErrorUnauthorized(responseError))
				return
		}
		if code == 500 {
				web.App.Logger.Error(err.Error())
				_ = render.Render(writer, request,
					errors.ErrorInternalServer(responseError))
				return
		}
		if !isAdmin {
				web.App.Logger.Error(responseError.Error())
				_ = render.Render(writer, request,
					errors.ErrorNotFound(responseError))
				return
		}
		next.ServeHTTP(writer, request)
	})
}