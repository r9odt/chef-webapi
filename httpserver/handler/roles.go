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

func rolesAPIHandler(writer http.ResponseWriter, request *http.Request) {
	responseError := fmt.Errorf("cannot get roles")
	setJSONHeader(writer)
	params := controller.GetURLAPIParameters(request)
	list, err := controller.GetRoles(params)
	if err != nil {
		web.App.Logger.Errorf(
			"rolesAPIHandler [controller.GetURLAPIParameters]: %s",
			err.Error())
		_ = render.Render(writer, request,
			errors.ErrorInternalServer(responseError))
		return
	}

	length := int64(len(list))
	setXTotalCountHeader(writer, length)

	var response APIResponse
	start, end := getStartAndEndFromParams(params, length)

	var roleNames = make([]string, 0)
	slice := list[start:end]
	for i := range slice {
		roleNames = append(roleNames, slice[i].ID)
	}
	nodesData := controller.GetRolesData(
		controller.ConstructSearchORFilter("role", roleNames), roleNames)
	for i := range slice {
		slice[i].Data = nodesData[slice[i].ID]
	}

	response, err = json.Marshal(list[start:end])
	if err != nil {
		web.App.Logger.Errorf(
			"rolesAPIHandler [json.Marshal]: %s",
			err.Error())
		_ = render.Render(writer, request,
			errors.ErrorInternalServer(responseError))
		return
	}

	fmt.Fprintf(writer, "%s", response)
}

func rolesGetNodesByRoleNameAPIHandler(writer http.ResponseWriter, request *http.Request) {
	setJSONHeader(writer)
	roleName := middleware.GetID(request)

	entry, err := controller.GetCompleteDeployersForResourceByName("roles",
		roleName)
	if err != nil {
		web.App.Logger.Errorf(
			"rolesGetNodesByRoleNameAPIHandler [controller.GetCompleteDeployersForResourceByName]: %s",
			err.Error())
		return
	}
	var roleNames = []string{roleName}
	nodesData := controller.GetRolesData(
		controller.ConstructSearchORFilter("role", roleNames), roleNames)
	resultList := &controller.Role{
		ID:   roleName,
		Data: nodesData[roleName],
		Date: entry.Date,
	}

	var response APIResponse
	response, err = json.Marshal(resultList)
	if err != nil {
		web.App.Logger.Errorf(
			"rolesGetNodesByRoleNameAPIHandler [json.Marshal]: %s",
			err.Error())
		_ = render.Render(writer, request,
			errors.ErrorInternalServer(err))
		return
	}

	fmt.Fprintf(writer, "%s\n", response)
}

func rolesCreateTaskForResourceAPIHandler(writer http.ResponseWriter, request *http.Request) {
	setJSONHeader(writer)
	role, err := controller.CreateTaskForRole(request)
	if err != nil {
		web.App.Logger.Errorf(
			"rolesCreateTaskForResourceAPIHandler [ontroller.CreateTaskForRole]: %s",
			err.Error())
		_ = render.Render(writer, request,
			errors.ErrorInternalServer(err))
		return
	}
	var response APIResponse
	response, err = json.Marshal(role)
	if err != nil {
		web.App.Logger.Errorf(
			"rolesCreateTaskForResourceAPIHandler [json.Marshal]: %s",
			err.Error())
		_ = render.Render(writer, request,
			errors.ErrorInternalServer(err))
		return
	}
	fmt.Fprintf(writer, "%s\n", response)
}
