package controller

import (
	"fmt"
	"strings"

	web "github.com/JIexa24/chef-webapi"

	mergeSort "github.com/JIexa24/go-mergeSort"
	"github.com/go-chef/chef"
)

// Cookbook is a structure describing cookbook.
type Cookbook struct {
	ID   string                `json:"id"`
	Meta chef.CookbookVersions `json:"meta"`
}

// GetCookbooks returns cookbook list on server.
func GetCookbooks(parameters *APIParameters) ([]Cookbook,
	error) {
	sortField := parameters.Sort
	sortOrder := parameters.Order
	query := parameters.Q
	client := web.App.GetChefClientConfig()
	if client == nil {
		return make([]Cookbook, 0), fmt.Errorf("client not configured")
	}
	list, err := client.Cookbooks.List()
	if err != nil {
		return nil, err
	}

	var sliceToSort = make([]mergeSort.Interface, 0)
	for k, v := range list {
		if strings.Contains(k, query) {
			n := Cookbook{
				ID:   k,
				Meta: v,
			}
			sliceToSort = append(sliceToSort, n)
		}
	}

	compareIDFunction := func(a, b mergeSort.Interface) bool {
		return a.(Cookbook).ID < b.(Cookbook).ID
	}

	switch sortField {
	default:
		sliceToSort = mergeSort.Sort(sliceToSort, compareIDFunction, false)
	case "id":
		switch sortOrder {
		case "DESC":
			sliceToSort = mergeSort.Sort(sliceToSort, compareIDFunction, true)
		default:
			sliceToSort = mergeSort.Sort(sliceToSort, compareIDFunction, false)
		}
	}

	var result = make([]Cookbook, len(sliceToSort))
	for i, v := range sliceToSort {
		result[i] = v.(Cookbook)
	}

	return result, nil
}
