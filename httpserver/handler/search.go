package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	web "github.com/JIexa24/chef-webapi"

	"github.com/JIexa24/chef-webapi/httpserver/controller"
	"github.com/JIexa24/chef-webapi/httpserver/errors"
	"github.com/JIexa24/chef-webapi/httpserver/middleware"

	"github.com/go-chi/render"
)

func searchAPIHandler(writer http.ResponseWriter, request *http.Request) {
	setJSONHeader(writer)
	searchIndex := middleware.GetSearchIndex(request)
	searchParam := middleware.GetSearchQuery(request)
	p := controller.GetURLAPIParameters(request)
	switch searchIndex {
	case
		"node",
		"role",
		"client",
		"environment":
		list, err := controller.GetResourceBySearch(searchIndex, searchParam, p.IncludeAllIntoParseResource)
		if err != nil {
			web.App.Logger.Errorf(
				"searchAPIHandler [controller.GetResourceBySearch]: %s",
				err.Error())
			_ = render.Render(writer, request,
				errors.ErrorInternalServer(err))
			return
		}
		if len(list) <= 0 {
			_ = render.Render(writer, request,
				errors.ErrorNotFound(fmt.Errorf("no objects")))
			return
		}

		var response APIResponse
		response, err = json.Marshal(list)
		if err != nil {
			web.App.Logger.Errorf(
				"searchAPIHandler [json.Marshal]: %s",
				err.Error())
			_ = render.Render(writer, request,
				errors.ErrorInternalServer(err))
			return
		}
		fmt.Fprintf(writer, "%s\n", response)
		return
	}
	_ = render.Render(writer, request,
		errors.ErrorInvalidRequest(fmt.Errorf("not allowed searchIndex")))
}
