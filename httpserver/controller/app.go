package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	web "github.com/JIexa24/chef-webapi"

	"github.com/JIexa24/chef-webapi/database/interfaces"
	"github.com/JIexa24/chef-webapi/httpserver/middleware"

	mergeSort "github.com/JIexa24/go-mergeSort"
)

type moduleRequest struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	IsON bool   `json:"isON"`
}

// GetAllAppModules return modules information.
func GetAllAppModules(parameters *APIParameters) ([]interfaces.AppModuleEntry,
	error) {
	sortField := parameters.Sort
	sortOrder := parameters.Order
	query := parameters.Q
	idFilter := parameters.ID
	list, err := web.App.DB.GetAllAppModules()
	if err != nil {
		return nil, err
	}

	var sliceToSort = make([]mergeSort.Interface, 0)
	for i := range list {
		if strings.Contains(list[i].Name, query) {
			n := interfaces.CopyAppModuleData(list[i])
			contain := false
			for i := range idFilter {
				if idFilter[i] == n.ID {
					contain = true
					break
				}
			}
			if len(idFilter) <= 0 || contain {
				sliceToSort = append(sliceToSort, *n)
			}
		}
	}

	compareNameFunction := func(a, b mergeSort.Interface) bool {
		return a.(interfaces.AppModuleEntry).Name <
			b.(interfaces.AppModuleEntry).Name
	}

	switch sortField {
	default:
		sliceToSort = mergeSort.Sort(sliceToSort, compareNameFunction, false)
	case "username":
		switch sortOrder {
		case "DESC":
			sliceToSort = mergeSort.Sort(sliceToSort, compareNameFunction, true)
		default:
			sliceToSort = mergeSort.Sort(sliceToSort, compareNameFunction, false)
		}
	}

	var result []interfaces.AppModuleEntry = make([]interfaces.AppModuleEntry,
		len(sliceToSort))
	for i, v := range sliceToSort {
		result[i] = v.(interfaces.AppModuleEntry)
	}

	return result, nil
}

// GetAppModuleByID return struct with module information.
// Information get by id.
func GetAppModuleByID(
	r *http.Request) (*interfaces.AppModuleEntry, error) {
		id := middleware.GetID(r)
		m, err := web.App.DB.GetAppModuleByID(id)
		if err != nil {
			return nil, err
		}
		if m == nil {
			return nil, fmt.Errorf("module not found")
		}
	
		return m, nil
}

// UpdateAppModuleByID updates module information.
// Information update by id
func UpdateAppModuleByID(
	r *http.Request) error {
	id := middleware.GetID(r)
	m, err := web.App.DB.GetAppModuleByID(id)
	if err != nil {
		return err
	}
	if m == nil {
		return fmt.Errorf("module not found")
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	updatedModule := &moduleRequest{}
	err = json.Unmarshal(body, updatedModule)
	if err != nil {
		return err
	}
	ImplementToAppModuleInfo(m, updatedModule)
	return web.App.DB.UpdateAppModuleByID(m.ID, m)
}

// ImplementToAppModuleInfo implement some updated 
// fields into AppModule's entry.
func ImplementToAppModuleInfo(
	module *interfaces.AppModuleEntry,
	updatedModule *moduleRequest) {
	if updatedModule == nil || module == nil {
		return
	}
  module.IsON = updatedModule.IsON
}
