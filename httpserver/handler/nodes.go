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

func nodesAPIHandler(writer http.ResponseWriter, request *http.Request) {
	setJSONHeader(writer)
	responseError := fmt.Errorf("cannot get nodes")
	params := controller.GetURLAPIParameters(request)
	list, err := controller.GetNodes(params)
	if err != nil {
		web.App.Logger.Errorf(
			"nodesAPIHandler [controller.GetNodes]: %s",
			err.Error())
		_ = render.Render(writer, request,
			errors.ErrorInternalServer(responseError))
		return
	}

	length := int64(len(list))
	setXTotalCountHeader(writer, length)

	var response APIResponse
	start, end := getStartAndEndFromParams(params, length)

	var nodeNames = make([]string, 0)
	slice := list[start:end]
	for i := range slice {
		nodeNames = append(nodeNames, slice[i].ID)
	}

	nodesData := controller.GetNodesData(
		controller.ConstructSearchORFilter("name", nodeNames))
	for i := range slice {
		slice[i].Data = nodesData[slice[i].ID]
	}

	response, err = json.Marshal(slice)
	if err != nil {
		web.App.Logger.Errorf(
			"nodesAPIHandler [json.Marshal]: %s",
			err.Error())
		_ = render.Render(writer, request,
			errors.ErrorInternalServer(responseError))
		return
	}

	fmt.Fprintf(writer, "%s", response)
}

func nodesGetNodesByNodeNameAPIHandler(writer http.ResponseWriter, request *http.Request) {
	setJSONHeader(writer)
	nodeName := middleware.GetID(request)
	entry, err := controller.GetCompleteDeployersForResourceByName("nodes",
		nodeName)
	if err != nil {
		web.App.Logger.Errorf(
			"nodesGetNodesByNodeNameAPIHandler [controller.GetCompleteDeployersForResourceByName]: %s",
			err.Error())
		return
	}

	nodesData := controller.GetNodesData(
		controller.ConstructSearchORFilter("name", []string{nodeName}))
	resultList := &controller.Node{
		ID:   nodeName,
		Data: nodesData[nodeName],
		Date: entry.Date,
	}

	var response APIResponse
	response, err = json.Marshal(resultList)
	if err != nil {
		web.App.Logger.Errorf(
			"nodesGetNodesByNodeNameAPIHandler [json.Marshal]: %s",
			err.Error())
		_ = render.Render(writer, request,
			errors.ErrorInternalServer(err))
		return
	}
	fmt.Fprintf(writer, "%s\n", response)
}

func nodesCreateTaskForResourceAPIHandler(writer http.ResponseWriter, request *http.Request) {
	setJSONHeader(writer)
	node, err := controller.CreateTaskForNode(request)
	if err != nil {
		web.App.Logger.Errorf(
			"nodesCreateTaskForResourceAPIHandler [controller.CreateTaskForNode]: %s",
			err.Error())
		_ = render.Render(writer, request,
			errors.ErrorInternalServer(err))
		return
	}
	var response APIResponse
	response, err = json.Marshal(node)
	if err != nil {
		web.App.Logger.Errorf(
			"nodesCreateTaskForResourceAPIHandler [json.Marshal]: %s",
			err.Error())
		_ = render.Render(writer, request,
			errors.ErrorInternalServer(err))
		return
	}
	fmt.Fprintf(writer, "%s\n", response)
}
