package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	web "github.com/r9odt/chef-webapi"

	"github.com/go-chi/render"

	"github.com/r9odt/chef-webapi/httpserver/controller"
	"github.com/r9odt/chef-webapi/httpserver/errors"
)

func authenticationAPIHandler(writer http.ResponseWriter,
	request *http.Request) {
	responseError := fmt.Errorf("authentication error")
	setJSONHeader(writer)
	if controller.IsUserBlocked(writer, request) {
		return
	}
	result, err := controller.Authenticate(request)
	if err != nil {
		web.App.Logger.Errorf(
			"authenticationAPIHandler [controller.Authenticate]: %s", err.Error())
		_ = render.Render(writer, request,
			errors.ErrorUnauthorized(responseError))
		return
	}

	if result == nil {
		_ = render.Render(writer, request,
			errors.ErrorUnauthorized(responseError))
		return
	}

	response, err := json.Marshal(result)
	if err != nil {
		web.App.Logger.Errorf("authenticationAPIHandler [json.Marshal]: %s",
			err.Error())
		_ = render.Render(writer, request,
			errors.ErrorUnauthorized(responseError))
		return
	}

	fmt.Fprintf(writer, "%s", response)
}

func authenticationLogoutAPIHandler(writer http.ResponseWriter,
	request *http.Request) {
	responseError := fmt.Errorf("logout error")
	setJSONHeader(writer)

	err := controller.Logout(request)
	if err != nil {
		web.App.Logger.Errorf(
			"authenticationLogoutAPIHandler [controller.Logout]: %s", err.Error())
		_ = render.Render(writer, request,
			errors.ErrorInternalServer(responseError))
		return
	}

	fmt.Fprint(writer, emptyJSON())
}

func authenticationPingAPIHandler(writer http.ResponseWriter,
	request *http.Request) {
	responseError := fmt.Errorf("cannot check session")
	setJSONHeader(writer)
	if controller.IsUserBlocked(writer, request) {
		return
	}

	result, err := controller.CheckSession(request)
	if err != nil {
		web.App.Logger.Errorf(
			"authenticationPingAPIHandler [controller.CheckSession]: %s",
			err.Error())
		_ = render.Render(writer, request,
			errors.ErrorInternalServer(responseError))
		return
	}
	response, err := json.Marshal(result)
	if err != nil {
		web.App.Logger.Errorf("authenticationPingAPIHandler [json.Marshal]: %s",
			err.Error())
		_ = render.Render(writer, request,
			errors.ErrorInternalServer(responseError))
		return
	}

	if !result.Authenticate {
		_ = render.Render(writer, request,
			errors.ErrorUnauthorized(responseError))
		return
	}
	fmt.Fprintf(writer, "%s", response)
}

func authenticationGetCurrentUserAPIHandler(writer http.ResponseWriter,
	request *http.Request) {
	responseError := fmt.Errorf("cannot get user")
	setJSONHeader(writer)
	if controller.IsUserBlocked(writer, request) {
		return
	}

	user, err := controller.GetUserBySession(request)
	if err != nil {
		web.App.Logger.Errorf(
			"authenticationGetCurrentUserAPIHandler [controller.GetUserBySession]: %s",
			err.Error())
		_ = render.Render(writer, request,
			errors.ErrorInternalServer(responseError))
	}

	if user == nil {
		_ = render.Render(writer, request,
			errors.ErrorUnauthorized(responseError))
		return
	}
	result := controller.ExtractAuthenticationUserInfo(user)
	response, err := json.Marshal(result)
	if err != nil {
		web.App.Logger.Errorf(
			"authenticationGetCurrentUserAPIHandler [json.Marshal]: %s",
			err.Error())
		_ = render.Render(writer, request,
			errors.ErrorInternalServer(responseError))
		return
	}

	fmt.Fprintf(writer, "%s", response)
}

func authenticationGetCurrentUserPermissionsAPIHandler(
	writer http.ResponseWriter,
	request *http.Request) {
	responseError := fmt.Errorf("cannot get user")
	setJSONHeader(writer)
	if controller.IsUserBlocked(writer, request) {
		_ = render.Render(writer, request,
			errors.ErrorUnauthorized(responseError))
		return
	}

	user, err := controller.GetUserBySession(request)
	if err != nil {
		web.App.Logger.Errorf(
			"authenticationGetCurrentUserAPIHandler [controller.GetUserBySession]: %s",
			err.Error())
		_ = render.Render(writer, request,
			errors.ErrorInternalServer(responseError))
		return
	}

	if user == nil {
		_ = render.Render(writer, request,
			errors.ErrorUnauthorized(responseError))
		return
	}
	result := controller.ExtractAuthenticationUserPermissions(user)
	response, err := json.Marshal(result)
	if err != nil {
		web.App.Logger.Errorf(
			"authenticationGetCurrentUserAPIHandler [json.Marshal]: %s",
			err.Error())
		_ = render.Render(writer, request,
			errors.ErrorInternalServer(responseError))
		return
	}

	fmt.Fprintf(writer, "%s", response)
}
