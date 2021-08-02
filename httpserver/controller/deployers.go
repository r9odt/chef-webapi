package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	web "github.com/JIexa24/chef-webapi"

	"github.com/JIexa24/chef-webapi/database/interfaces"
	"github.com/JIexa24/chef-webapi/httpserver/middleware"
)

// GetAllTasks return tasks information.
func GetAllTasks(parameters *APIParameters) ([]interfaces.TaskEntry, error) {
	tasks, err := web.App.DB.GetAllTasks()
	for i, v := range tasks {
		tasks[i].Log = getLogForDeployer(v.ID, v.Resource, v.Name)
	}
	return tasks, err
}

type deployerRequest struct {
	Name             string `json:"id"`
	OnlyResource     bool   `json:"onlyResource"`
	Resources        string `json:"resources"`
	SelectedResource bool   `json:"selectedResource"`
}

// CreateTask creates a task.
func CreateTask(r *http.Request) (*interfaces.TaskEntry, error) {
	resource := middleware.GetID(r)
	user, err := GetUserBySession(r)
	if user == nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	var req = &deployerRequest{}

	err = json.Unmarshal(body, req)
	if err != nil {
		return nil, err
	}
	return CreateTaskFromRequest(resource, user.ID, req)
}

// CreateTaskFromRequest creating task that spec in fields of deployerRequest.
// Also check for task existing by resource and name.
func CreateTaskFromRequest(resource string, userID string,
	req *deployerRequest) (*interfaces.TaskEntry, error) {
	taskExist := web.App.DB.СheckIfTaskAlreadyCreate(resource, req.Name)
	if taskExist {
		return nil, fmt.Errorf("task for %s: %s already exist", resource, req.Name)
	}

	switch resource {
	case "nodes":
		// We cannot deploy only resource if resource is a node
		req.OnlyResource = false
	}

	return web.App.DB.CreateTask(resource, req.Name,
		req.Resources, userID, req.OnlyResource, req.SelectedResource)
}

// GetCompleteDeployersForResourceByName returns the timestamp
// of the last deployment task for the specified resource.
func GetCompleteDeployersForResourceByName(resource string, name string) (*interfaces.TaskEntry, error) {
	noDeployedMessage := "No deployed yet"
	errorFetchMessage := "Unable to fecth information"
	task := interfaces.NewEmptyTask()
	entry, err := web.App.DB.GetLastCompleteTaskByResourceAndName(resource, name)
	if err != nil {
		task.Date = errorFetchMessage
		task.InitiatorID = ""
		task.Status = errorFetchMessage
		return task, err
	}
	if entry != nil {
		task = entry
	} else {
		task.Date = noDeployedMessage
		task.InitiatorID = ""
		task.Status = noDeployedMessage
	}

	return task, nil
}

// GetDeployerByID return task with the specified ID.
func GetDeployerByID(r *http.Request) (*interfaces.TaskEntry, error) {
	id := middleware.GetID(r)
	if id == "" {
		return nil, fmt.Errorf("id not found")
	}

	entry, err := web.App.DB.GetTaskByID(id)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}
	response := entry
	response.Log = getLogForDeployer(entry.ID, entry.Resource, entry.Name)
	return response, nil
}

func getLogForDeployer(id string, resource, name string) string {
	resourceFile := fmt.Sprintf("worker-%s-%s-%s.log",
		id, resource, name)
	var b []byte
	var err error
	if !web.App.DB.CheckFile(resourceFile) {
		b, err = os.ReadFile(filepath.Join(web.App.WorkerDirectory, resourceFile))
		if err != nil {
			return "Cannot read file"
		}
	} else {
		b = web.App.DB.DownloadFile(resourceFile)
	}
	return string(b)
}
