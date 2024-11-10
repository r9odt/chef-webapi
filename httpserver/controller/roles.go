package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	web "github.com/r9odt/chef-webapi"

	"github.com/r9odt/chef-webapi/database/interfaces"
	"github.com/r9odt/chef-webapi/httpserver/middleware"

	mergeSort "github.com/r9odt/go-mergeSort"
)

// Role is a structure describing role.
type Role struct {
	ID   string                   `json:"id"`
	Date string                   `json:"date"`
	Data []map[string]interface{} `json:"data"`
}

// GetRoles returns role list on server.
func GetRoles(parameters *APIParameters) ([]Role, error) {
	sortField := parameters.Sort
	sortOrder := parameters.Order
	query := parameters.Q
	client := web.App.GetChefClientConfig()
	if client == nil {
		return make([]Role, 0), fmt.Errorf("client not configured")
	}
	list, err := client.Roles.List()
	if err != nil {
		return nil, err
	}

	var sliceToSort = make([]mergeSort.Interface, 0)
	for k := range *list {
		if strings.Contains(k, query) {
			entry, err := GetCompleteDeployersForResourceByName("roles", k)
			if err != nil {
				web.App.Logger.Errorf(
					"GetRoles [GetDeployersForResourceByNames]: %s",
					err.Error())
			}
			n := Role{
				ID:   k,
				Date: entry.Date,
				Data: nil,
			}
			sliceToSort = append(sliceToSort, n)
		}
	}

	compareIDFunction := func(a, b mergeSort.Interface) bool {
		return a.(Role).ID < b.(Role).ID
	}
	compareDateFunction := func(a, b mergeSort.Interface) bool {
		ta, _ := time.Parse(interfaces.TimeFormat, a.(Role).Date)
		tb, _ := time.Parse(interfaces.TimeFormat, b.(Role).Date)
		return ta.Unix() < tb.Unix()
	}

	sliceToSort = mergeSort.Sort(sliceToSort, compareIDFunction, false)
	switch sortField {
	// Case id is default because resource name uses as ID.
	// Chaos if resource was not sorted by ID.
	default:
		sliceToSort = mergeSort.Sort(sliceToSort, compareIDFunction, false)
	case "id":
		switch sortOrder {
		case "DESC":
			sliceToSort = mergeSort.Sort(sliceToSort, compareIDFunction, true)
		}
	case "date":
		switch sortOrder {
		case "DESC":
			sliceToSort = mergeSort.Sort(sliceToSort, compareDateFunction, true)
		default:
			sliceToSort = mergeSort.Sort(sliceToSort, compareDateFunction, false)
		}
	}

	var result = make([]Role, len(sliceToSort))
	for i, v := range sliceToSort {
		result[i] = v.(Role)
	}

	return result, nil
}

// GetRolesData return a map with node information for each role,
// specified by roles[]
func GetRolesData(Filter string,
	roles []string) map[string][]map[string]interface{} {
	list, err := GetResourceBySearch("node", Filter, false)
	if err != nil {
		web.App.Logger.Errorf(
			"GetRolesData [GetResourceBySearch]: %s",
			err.Error())
		return nil
	}

	result := make(map[string][]map[string]interface{})
	for _, v := range roles {
		for j := range list {
			index := v
			roleString := fmt.Sprintf("role[%s]", index)
			var interfaceArray = list[j]["run_list"].([]map[string]interface{})
			stringArray := make([]string, len(interfaceArray))
			for k := 0; k < len(interfaceArray); k++ {
				stringArray[k] = interfaceArray[k]["object"].(string)
			}

			if checkStringInArray(roleString, stringArray) {
				if result[index] == nil {
					result[index] = make([]map[string]interface{}, 0)
				}
				result[index] = append(result[index], list[j])
			}
		}
	}
	return result
}

func checkStringInArray(str string, array []string) bool {
	for i := 0; i < len(array); i++ {
		if array[i] == str {
			return true
		}
	}
	return false
}

// RoleTaskRequest for create task for role.
type RoleTaskRequest struct {
	Role
	OnlyResource bool `json:"onlyResource"`
}

// CreateTaskForRole creates task by spec request.
func CreateTaskForRole(r *http.Request) (*RoleTaskRequest, error) {
	roleName := middleware.GetID(r)
	resource := "roles"
	user, err := GetUserBySession(r)
	if user == nil {
		return nil, err
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	var req = &RoleTaskRequest{}

	err = json.Unmarshal(body, req)
	if err != nil {
		return nil, err
	}

	var resources = make([]string, 0)
	var resourceList = req.Data
	for i := range resourceList {
		isSelected := resourceList[i]["selected"].(bool)
		if isSelected {
			resources = append(resources, resourceList[i]["name"].(string))
		}
	}
	resourcesJSON, err := json.Marshal(resources)
	if err != nil {
		return nil, err
	}

	deployerRequest := &deployerRequest{
		Name:             roleName,
		OnlyResource:     req.OnlyResource,
		SelectedResource: true,
		Resources:        string(resourcesJSON),
	}

	_, err = CreateTaskFromRequest(resource, user.ID, deployerRequest)

	if err != nil {
		return nil, err
	}

	return req, nil
}
