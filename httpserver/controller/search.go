package controller

import (
	"fmt"

	web "github.com/r9odt/chef-webapi"

	"github.com/go-chef/chef"
	mergeSort "github.com/r9odt/go-mergeSort"
)

// GetResourceBySearch return list what contains list of Resource on the server
// with specified search.
func GetResourceBySearch(index, val string, dontCutParseResult bool) ([]map[string]interface{}, error) {
	if val == "*:*" {
		return nil, fmt.Errorf("not allowed pattern *:*")
	}
	if val == "" {
		return nil, nil
	}

	client := web.App.GetChefClientConfig()
	if client == nil {
		return make([]map[string]interface{}, 0),
			fmt.Errorf("client not configured")
	}
	query, err := client.Search.NewQuery(index, val)
	if err != nil {
		return nil, err
	}

	// Run the query
	list, err := query.Do(client)
	if err != nil {
		return nil, err
	}
	var result []map[string]interface{}
	// Processing search
	switch index {
	case "node":
		result = parseNodes(&list, dontCutParseResult)
	case "role":
		result = parseRoles(&list)
	case "client":
		result = parseClients(&list)
	case "environment":
		result = parseEnvs(&list)
	}
	return result, nil
}

func parseNodes(list *chef.SearchResult, dontCutParseResult bool) []map[string]interface{} {
	var sliceToSort = make([]mergeSort.Interface, 0)

	for _, v := range list.Rows {
		el := make(map[string]interface{})
		el["name"] = v.(map[string]interface{})["name"]
		el["selected"] = false
		runList := v.(map[string]interface{})["run_list"]
		resultRunList := make([]map[string]interface{}, 0)
		if runList != nil {
			runListArray := runList.([]interface{})
			for i := range runListArray {
				runListElement := make(map[string]interface{})
				runListElement["selected"] = false
				runListElement["object"] = runListArray[i].(string)
				resultRunList = append(resultRunList, runListElement)
			}
		}
		el["run_list"] = resultRunList
		el["ipaddress"] =
			v.(map[string]interface{})["automatic"].(map[string]interface{})["ipaddress"]
		el["fqdn"] = v.(map[string]interface{})["automatic"].(map[string]interface{})["fqdn"]
		el["ohai_time"] = v.(map[string]interface{})["automatic"].(map[string]interface{})["ohai_time"]
		sliceToSort = append(sliceToSort, el)
		if dontCutParseResult {
			el["data"] = v
		}
	}

	compareNameFunction := func(a, b mergeSort.Interface) bool {
		return a.(map[string]interface{})["name"].(string) <
			b.(map[string]interface{})["name"].(string)
	}
	sliceToSort = mergeSort.Sort(sliceToSort, compareNameFunction, false)

	result := make([]map[string]interface{}, len(sliceToSort))
	for i, v := range sliceToSort {
		result[i] = v.(map[string]interface{})
	}

	return result
}

func parseRoles(list *chef.SearchResult) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)
	for _, v := range list.Rows {
		el := make(map[string]interface{})
		el["res"] = v.(map[string]interface{})
		result = append(result, el)
	}
	return result
}

func parseClients(list *chef.SearchResult) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)
	for _, v := range list.Rows {
		el := make(map[string]interface{})
		el["res"] = v.(map[string]interface{})
		result = append(result, el)
	}
	return result
}

func parseEnvs(list *chef.SearchResult) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)
	for _, v := range list.Rows {
		el := make(map[string]interface{})
		el["res"] = v.(map[string]interface{})
		result = append(result, el)
	}
	return result
}

// ConstructSearchORFilter constructing string like
// 'names[0] OR names[1] OR ... OR names[n]'
func ConstructSearchORFilter(sign string, names []string) string {
	var result string
	if sign == "" {
		return result
	}
	var length = len(names)
	var orstring = " OR "
	for i, v := range names {
		result = result + fmt.Sprintf("%s:%s", sign, v)
		if i < length-1 {
			result = result + orstring
		}
	}
	return result
}
