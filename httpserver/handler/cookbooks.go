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

func cookbooksAPIHandler(writer http.ResponseWriter, request *http.Request) {
	setJSONHeader(writer)
	responseError := fmt.Errorf("cannot get cookbooks")
	params := controller.GetURLAPIParameters(request)
	list, err := controller.GetCookbooks(params)
	if err != nil {
		web.App.Logger.Errorf(
			"cookbooksAPIHandler [controller.GetCookbooks]: %s",
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
			"cookbooksAPIHandler [json.Marshal]: %s",
			err.Error())
		_ = render.Render(writer, request,
			errors.ErrorInternalServer(responseError))
		return
	}

	fmt.Fprintf(writer, "%s", response)
}
