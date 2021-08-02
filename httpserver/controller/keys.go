package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	web "github.com/JIexa24/chef-webapi"

	"github.com/JIexa24/chef-webapi/httpserver/middleware"

	mergeSort "github.com/JIexa24/go-mergeSort"
)

// KeyEntry is a struct describes application key.
type KeyEntry struct {
	Name  string `json:"id"`
	Value string `json:"value"`
}

// GetAllAppKeys return all application keys.
func GetAllAppKeys(parameters *APIParameters) ([]KeyEntry, error) {
	sortField := parameters.Sort
	sortOrder := parameters.Order
	query := parameters.Q
	idFilter := parameters.ID

	list := make([]KeyEntry, 2)
	list[0].Name = "SSH"
	key, _ := os.ReadFile(web.App.SSHKeyPath)
	list[0].Value = string(key)
	list[1].Name = "Chef"
	key, _ = os.ReadFile(web.App.ChefKeyPath)
	list[1].Value = string(key)

	var sliceToSort = make([]mergeSort.Interface, 0)
	for i := range list {
		if strings.Contains(list[i].Name, query) {
			n := KeyEntry{
				Name:  list[i].Name,
				Value: list[i].Value,
			}
			contain := false
			for i := range idFilter {
				if idFilter[i] == n.Name {
					contain = true
					break
				}
			}
			if len(idFilter) <= 0 || contain {
				sliceToSort = append(sliceToSort, n)
			}
		}
	}

	compareNameFunction := func(a, b mergeSort.Interface) bool {
		return a.(KeyEntry).Name <
			b.(KeyEntry).Name
	}

	switch sortField {
	default:
		sliceToSort = mergeSort.Sort(sliceToSort, compareNameFunction, false)
	case "id":
		switch sortOrder {
		case "DESC":
			sliceToSort = mergeSort.Sort(sliceToSort, compareNameFunction, true)
		default:
			sliceToSort = mergeSort.Sort(sliceToSort, compareNameFunction, false)
		}
	}

	var result = make([]KeyEntry,
		len(sliceToSort))
	for i, v := range sliceToSort {
		result[i] = v.(KeyEntry)
	}

	return result, nil
}

// GetAppKeyByID return struct with key information.
// Information get by id.
func GetAppKeyByID(r *http.Request) (*KeyEntry, error) {
	id := middleware.GetID(r)
	switch id {
	case "SSH":
		key, _ := os.ReadFile(web.App.SSHKeyPath)
		return &KeyEntry{
			Name:  id,
			Value: string(key),
		}, nil
	case "Chef":
		key, _ := os.ReadFile(web.App.ChefKeyPath)
		return &KeyEntry{
			Name:  id,
			Value: string(key),
		}, nil
	}
	return nil, fmt.Errorf("module not found")
}

// SetAppKeyByID updates key information.
// Information update by id.
func SetAppKeyByID(r *http.Request) error {
	id := middleware.GetID(r)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	k := &KeyEntry{}
	err = json.Unmarshal(body, k)
	if err != nil {
		return err
	}

	switch id {
	case "SSH":
		_ = os.WriteFile(web.App.SSHKeyPath, []byte(k.Value), 0600)
	case "Chef":
		_ = os.WriteFile(web.App.ChefKeyPath, []byte(k.Value), 0600)
	}

	web.App.ReloadChannel <- struct{}{}

	return nil
}
