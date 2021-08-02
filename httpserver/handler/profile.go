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

func profileGetUserProfileByIDAPIHandler(writer http.ResponseWriter,
	request *http.Request) {
	responseError := fmt.Errorf("cannot get profile")
	setJSONHeader(writer)
	profile, err := controller.GetUserProfileByID(request)
	if err != nil {
		web.App.Logger.Errorf(
			"profileGetUserProfileByIDAPIHandler [controller.GetUserProfileByID]: %s",
			err.Error())
		_ = render.Render(writer, request,
			errors.ErrorInternalServer(responseError))
		return
	}

	var response APIResponse
	response, err = json.Marshal(profile)
	if err != nil {
		web.App.Logger.Errorf(
			"profileGetUserProfileByIDAPIHandler [json.Marshal]: %s",
			err.Error())
		_ = render.Render(writer, request,
			errors.ErrorInternalServer(responseError))
		return
	}

	fmt.Fprintf(writer, "%s", response)
}

func profileGetCurrentUserAPIHandler(writer http.ResponseWriter,
	request *http.Request) {
	responseError := fmt.Errorf("cannot get profile")
	setJSONHeader(writer)
	if controller.IsUserBlocked(writer, request) {
		return
	}

	user, err := controller.GetUserBySession(request)
	if err != nil {
		web.App.Logger.Errorf(
			"profileGetCurrentUserAPIHandler [controller.GetUserBySession]: %s",
			err.Error())
		_ = render.Render(writer, request,
			errors.ErrorInternalServer(responseError))
	}

	if user == nil {
		_ = render.Render(writer, request,
			errors.ErrorUnauthorized(responseError))
		return
	}

	result := controller.ExtractProfileUserInfo(user)
	result.ID = "edit"
	result.UserID = user.ID
	result.Username = user.Username
	response, err := json.Marshal(result)
	if err != nil {
		web.App.Logger.Errorf(
			"profileGetCurrentUserAPIHandler [json.Marshal]: %s",
			err.Error())
		_ = render.Render(writer, request,
			errors.ErrorInternalServer(responseError))
		return
	}

	fmt.Fprintf(writer, "%s", response)
}

func profileUpdateCurrentUserAPIHandler(writer http.ResponseWriter,
	request *http.Request) {
	responseError := fmt.Errorf("cannot update profile")
	setJSONHeader(writer)
	err := controller.UpdateCurrentUserProfile(request)
	if err != nil {
		web.App.Logger.Errorf(
			"usersUpdateUserByIDAPIHandler [controller.UpdateUserByID]: %s",
			err.Error())
		_ = render.Render(writer, request,
			errors.ErrorInternalServer(responseError))
		return
	}
	profileGetCurrentUserAPIHandler(writer, request)
}

func profileGetUserProfilesAPIHandler(writer http.ResponseWriter,
	request *http.Request) {
	responseError := fmt.Errorf("cannot get profile")
	setJSONHeader(writer)

	params := controller.GetURLAPIParameters(request)
	profiles, err := controller.GetUserProfiles(params)
	if err != nil {
		web.App.Logger.Errorf(
			"profileGetUserProfilesAPIHandler [controller.GetUserProfiles]: %s",
			err.Error())
		_ = render.Render(writer, request,
			errors.ErrorInternalServer(responseError))
		return
	}

	length := int64(len(profiles))
	setXTotalCountHeader(writer, length)

	var response APIResponse
	response, err = json.Marshal(profiles)
	if err != nil {
		web.App.Logger.Errorf(
			"profileGetUserProfilesAPIHandler [json.Marshal]: %s",
			err.Error())
		_ = render.Render(writer, request,
			errors.ErrorInternalServer(responseError))
		return
	}

	fmt.Fprintf(writer, "%s", response)
}
