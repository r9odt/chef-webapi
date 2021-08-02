package controller

import (
	"fmt"
	"net/http"

	"github.com/JIexa24/chef-webapi/httpserver/errors"

	"github.com/go-chi/render"
)

// IsUserBlocked checks if the user is locked.
func IsUserBlocked(w http.ResponseWriter, r *http.Request) bool {
	user, err := GetUserBySession(r)
	if err != nil || user == nil {
		return false
	}
	if GetUserIsBlocked(r) {
		_ = render.Render(w, r,
			errors.ErrorUnauthorized(fmt.Errorf("Error")))
		return true
	}
	return false
}
