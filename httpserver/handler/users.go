package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	web "github.com/JIexa24/chef-webapi"

	"github.com/JIexa24/chef-webapi/httpserver/controller"
	"github.com/JIexa24/chef-webapi/httpserver/errors"

	"github.com/go-chi/render"
)

func usersDeleteUserByIDAPIHandler(writer http.ResponseWriter,
	request *http.Request) {
	responseError := fmt.Errorf("cannot delete user")
	setJSONHeader(writer)
	err := controller.DeleteUserByID(request)
	if err != nil {
		web.App.Logger.Errorf(
			"usersDeleteUserByIDAPIHandler [controller.DeleteUserByID]: %s",
			err.Error())
		_ = render.Render(writer, request,
			errors.ErrorInternalServer(responseError))
		return
	}

	fmt.Fprint(writer, emptyJSON())
}

func usersUpdateUserByIDAPIHandler(writer http.ResponseWriter,
	request *http.Request) {
	responseError := fmt.Errorf("cannot update user")
	setJSONHeader(writer)
	err := controller.UpdateUserByID(request)
	if err != nil {
		web.App.Logger.Errorf(
			"usersUpdateUserByIDAPIHandler [controller.UpdateUserByID]: %s",
			err.Error())
		_ = render.Render(writer, request,
			errors.ErrorInternalServer(responseError))
		return
	}
	usersGetUserByIDAPIHandler(writer, request)
}

func usersGetUserByIDAPIHandler(writer http.ResponseWriter,
	request *http.Request) {
	responseError := fmt.Errorf("cannot get user")
	setJSONHeader(writer)

	user, err := controller.GetUserByID(request)
	if err != nil {
		web.App.Logger.Errorf(
			"usersGetUserByIDAPIHandler [controller.GetUserByID]: %s",
			err.Error())
		_ = render.Render(writer, request,
			errors.ErrorNotFound(responseError))
		return
	}

	response, err := json.Marshal(user)
	if err != nil {
		web.App.Logger.Errorf(
			"usersGetUserByIDAPIHandler [json.Marshal]: %s",
			err.Error())
		_ = render.Render(writer, request,
			errors.ErrorInternalServer(responseError))
		return
	}

	fmt.Fprintf(writer, "%s", response)
}

func usersCreateUserAPIHandler(writer http.ResponseWriter,
	request *http.Request) {
	responseError := fmt.Errorf("cannot create user")
	setJSONHeader(writer)
	result, err := controller.CreateUser(request)
	if err != nil {
		_ = render.Render(writer, request,
			errors.ErrorInternalServer(err))
		return
	}

	response, err := json.Marshal(result)
	if err != nil {
		web.App.Logger.Errorf(
			"usersCreateUserAPIHandler [json.Marshal]: %s",
			err.Error())
		_ = render.Render(writer, request,
			errors.ErrorInternalServer(responseError))
		return
	}

	fmt.Fprintf(writer, "%s", response)
}

func usersGetUsersAPIHandler(writer http.ResponseWriter,
	request *http.Request) {
	responseError := fmt.Errorf("cannot get users")
	setJSONHeader(writer)
	params := controller.GetURLAPIParameters(request)

	list, err := controller.GetAllUsers(params)
	if err != nil {
		web.App.Logger.Errorf(
			"usersGetUsersAPIHandler [controller.GetAllUsers]: %s",
			err.Error())
		_ = render.Render(writer, request,
			errors.ErrorInternalServer(responseError))
		return
	}

	length := int64(len(list))
	setXTotalCountHeader(writer, length)

	var response APIResponse
	start, end := getStartAndEndFromParams(params, length)

	response, err = json.Marshal(list[start:end])
	if err != nil {
		web.App.Logger.Errorf(
			"usersGetUsersAPIHandler [json.Marshal]: %s",
			err.Error())
		_ = render.Render(writer, request,
			errors.ErrorInternalServer(responseError))
		return
	}
	fmt.Fprintf(writer, "%s", response)
}
