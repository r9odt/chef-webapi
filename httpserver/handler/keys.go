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

func keysGetKeysAPIHandler(writer http.ResponseWriter,
	request *http.Request) {
	responseError := fmt.Errorf("cannot get keys")
	setJSONHeader(writer)
	params := controller.GetURLAPIParameters(request)

	list, err := controller.GetAllAppKeys(params)
	if err != nil {
		web.App.Logger.Errorf(
			"keysGetKeysAPIHandler [controller.GetAllAppKeys]: %s",
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
			"keysGetKeysAPIHandler [json.Marshal]: %s",
			err.Error())
		_ = render.Render(writer, request,
			errors.ErrorInternalServer(responseError))
		return
	}
	fmt.Fprintf(writer, "%s", response)
}

func keysGetKeyByIDAPIHandler(writer http.ResponseWriter,
	request *http.Request) {
	responseError := fmt.Errorf("cannot get keys")
	setJSONHeader(writer)

	m, err := controller.GetAppKeyByID(request)
	if err != nil {
		web.App.Logger.Errorf(
			"keysGetKeyByIDAPIHandler [controller.GetAppKeyByID]: %s",
			err.Error())
		_ = render.Render(writer, request,
			errors.ErrorInternalServer(responseError))
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

func keysUpdateKeyByIDAPIHandler(writer http.ResponseWriter,
	request *http.Request) {
	responseError := fmt.Errorf("cannot update key")
	setJSONHeader(writer)
	err := controller.SetAppKeyByID(request)
	if err != nil {
		web.App.Logger.Errorf(
			"keysUpdateKeyByIDAPIHandler [controller.SetAppKeyByID]: %s",
			err.Error())
		_ = render.Render(writer, request,
			errors.ErrorInternalServer(responseError))
		return
	}
	keysGetKeyByIDAPIHandler(writer, request)
}
