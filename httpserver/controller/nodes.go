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

// Node is a structure describing node.
type Node struct {
	ID   string                   `json:"id"`
	Date string                   `json:"date"`
	Data []map[string]interface{} `json:"data"`
}

// GetNodes returns node list on server.
func GetNodes(parameters *APIParameters) ([]Node, error) {
	sortField := parameters.Sort
	sortOrder := parameters.Order
	query := parameters.Q
	idFilter := parameters.ID
	client := web.App.GetChefClientConfig()
	if client == nil {
		return make([]Node, 0), fmt.Errorf("client not configured")
	}
	list, err := client.Nodes.List()
	if err != nil {
		return nil, err
	}

	var sliceToSort = make([]mergeSort.Interface, 0)
	for k := range list {
		if strings.Contains(k, query) {
			entry, err := GetCompleteDeployersForResourceByName("nodes", k)
			if err != nil {
				web.App.Logger.Errorf(
					"GetRoles [GetDeployersForResourceByNames]: %s",
					err.Error())
			}
			n := Node{
				ID:   k,
				Date: entry.Date,
				Data: nil,
			}
			contain := false
			for i := range idFilter {
				if idFilter[i] == n.ID {
					contain = true
					break
				}
			}
			if len(idFilter) <= 0 || contain {
				sliceToSort = append(sliceToSort, n)
			}
		}
	}

	compareIDFunction := func(a, b mergeSort.Interface) bool {
		return a.(Node).ID < b.(Node).ID
	}
	compareDateFunction := func(a, b mergeSort.Interface) bool {
		ta, _ := time.Parse(interfaces.TimeFormat, a.(Node).Date)
		tb, _ := time.Parse(interfaces.TimeFormat, b.(Node).Date)
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

	var result = make([]Node, len(sliceToSort))
	for i, v := range sliceToSort {
		result[i] = v.(Node)
	}

	return result, nil
}

// GetNodesData return a map with node information for each node,
// with filter.
func GetNodesData(filter string) map[string][]map[string]interface{} {
	list, err := GetResourceBySearch("node", filter, false)
	if err != nil {
		web.App.Logger.Errorf(
			"GetNodesData [GetResourceBySearch]: %s",
			err.Error())
		return nil
	}

	result := make(map[string][]map[string]interface{})
	for i, v := range list {
		index := v["name"].(string)
		if result[index] == nil {
			result[index] = make([]map[string]interface{}, 0)
		}
		result[index] = append(result[index], list[i])
	}
	return result
}

// NodeTaskRequest for create task for node.
type NodeTaskRequest struct {
	Node
	OnlyResource bool `json:"onlyResource"`
}

// CreateTaskForNode creates task by spec request.
func CreateTaskForNode(r *http.Request) (*NodeTaskRequest, error) {
	nodeName := middleware.GetID(r)
	resource := "nodes"
	user, err := GetUserBySession(r)
	if user == nil {
		return nil, err
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	var req = &NodeTaskRequest{}

	err = json.Unmarshal(body, req)
	if err != nil {
		return nil, err
	}

	var resources = make([]string, 0)
	var nodesList = req.Data
	var resourceList = nodesList[0]["run_list"].([]interface{})
	for i := range resourceList {
		res := resourceList[i].(map[string]interface{})
		isSelected := res["selected"].(bool)
		if isSelected {
			resources = append(resources, res["object"].(string))
		}
	}
	resourcesJSON, err := json.Marshal(resources)
	if err != nil {
		return nil, err
	}

	deployerRequest := &deployerRequest{
		Name:             nodeName,
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
