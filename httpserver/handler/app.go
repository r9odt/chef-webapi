package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	web "github.com/r9odt/chef-webapi"

	"github.com/r9odt/chef-webapi/httpserver/controller"
	"github.com/r9odt/chef-webapi/httpserver/errors"

	"github.com/go-chi/render"
)

func appGetModulesAPIHandler(
	writer http.ResponseWriter,
	request *http.Request) {
	responseError := fmt.Errorf("cannot get modules")
	setJSONHeader(writer)
	params := controller.GetURLAPIParameters(request)

	list, err := controller.GetAllAppModules(params)
	if err != nil {
		web.App.Logger.Errorf(
			"appGetModulesAPIHandler [controller.GetAllAppModules]: %s",
			err.Error())
		_ = render.Render(writer, request,
			errors.ErrorInternalServer(responseError))
		return
	}

	length := int64(len(list))
	setXTotalCountHeader(writer, length)

	var response APIResponse
	start := params.Start
	if start < 0 {
		start = 0
	}
	end := params.End
	if end >= length {
		end = length
	}

	response, err = json.Marshal(list[start:end])
	if err != nil {
		web.App.Logger.Errorf(
			"appGetModulesAPIHandler [json.Marshal]: %s",
			err.Error())
		_ = render.Render(writer, request,
			errors.ErrorInternalServer(responseError))
		return
	}
	fmt.Fprintf(writer, "%s", response)
}

func appGetModuleByIDAPIHandler(
	writer http.ResponseWriter,
	request *http.Request) {
	responseError := fmt.Errorf("cannot get module")
	setJSONHeader(writer)

	m, err := controller.GetAppModuleByID(request)
	if err != nil {
		web.App.Logger.Errorf(
			"appGetModuleByIDAPIHandler [controller.GetAppModuleByID]: %s",
			err.Error())
		_ = render.Render(writer, request,
			errors.ErrorNotFound(responseError))
		return
	}

	response, err := json.Marshal(m)
	if err != nil {
		web.App.Logger.Errorf(
			"appGetModuleByIDAPIHandler [json.Marshal]: %s",
			err.Error())
		_ = render.Render(writer, request,
			errors.ErrorInternalServer(responseError))
		return
	}

	fmt.Fprintf(writer, "%s", response)
}

func appUpdateModuleByIDAPIHandler(
	writer http.ResponseWriter,
	request *http.Request) {
	responseError := fmt.Errorf("cannot update module")
	setJSONHeader(writer)
	err := controller.UpdateAppModuleByID(request)
	if err != nil {
		web.App.Logger.Errorf(
			"appUpdateModuleByIDAPIHandler [controller.UpdateAppModuleByID]: %s",
			err.Error())
		_ = render.Render(writer, request,
			errors.ErrorInternalServer(responseError))
		return
	}
	appGetModuleByIDAPIHandler(writer, request)
}
