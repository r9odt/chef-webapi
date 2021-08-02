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

func deployersAPIHandler(writer http.ResponseWriter, request *http.Request) {
	responseError := fmt.Errorf("cannot get tasks")
	setJSONHeader(writer)
	params := controller.GetURLAPIParameters(request)

	list, err := controller.GetAllTasks(params)
	if err != nil {
		web.App.Logger.Errorf(
			"deployersAPIHandler [controller.GetAllTasks]: %s",
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
			"deployersAPIHandler [json.Marshal]: %s",
			err.Error())
		_ = render.Render(writer, request,
			errors.ErrorInternalServer(responseError))
		return
	}
	fmt.Fprintf(writer, "%s", response)
}

func deployersCreateDeployerAPIHandler(writer http.ResponseWriter, request *http.Request) {
	setJSONHeader(writer)
	apiResponse, err := controller.CreateTask(request)
	if err != nil {
		web.App.Logger.Errorf(
			"deployersCreateDeployerAPIHandler [controller.CreateTask]: %s",
			err.Error())
		_ = render.Render(writer, request,
			errors.ErrorNotAcceptable(err))
		return
	}

	response, err := json.Marshal(apiResponse)
	if err != nil {
		web.App.Logger.Errorf(
			"deployersCreateDeployerAPIHandler [json.Marshal]: %s",
			err.Error())
		_ = render.Render(writer, request,
			errors.ErrorInternalServer(fmt.Errorf("cannot get task info")))
		return
	}

	fmt.Fprintf(writer, "%s", response)
}

func deployersGetDeployerByIDAPIHandler(writer http.ResponseWriter, request *http.Request) {
	setJSONHeader(writer)

	log, err := controller.GetDeployerByID(request)
	if err != nil {
		web.App.Logger.Errorf(
			"deployersGetDeployerLogAPIHandler [controller.GetDeployerLog]: %s",
			err.Error())
	}
	if log == nil {
		_ = render.Render(writer, request,
			errors.ErrorNotFound(fmt.Errorf("log not found")))
		return
	}
	response, err := json.Marshal(log)
	if err != nil {
		web.App.Logger.Errorf(
			"deployersGetDeployerLogAPIHandler [json.Marshal]: %s",
			err.Error())
		_ = render.Render(writer, request,
			errors.ErrorInternalServer(fmt.Errorf("cannot get task info")))
		return
	}
	fmt.Fprintf(writer, "%s", response)
}
