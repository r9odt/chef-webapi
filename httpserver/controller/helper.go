package controller

import (
	"net/http"
	"strconv"
)

// APIParameters contains all parameter, which can be processing ftom url query.
type APIParameters struct {
	Start                       int64
	End                         int64
	Sort                        string
	Order                       string
	Q                           string
	ID                          []string
	IncludeAllIntoParseResource bool
}

// GetURLAPIParameters return APIParameters from url args.
func GetURLAPIParameters(r *http.Request) *APIParameters {
	var params = &APIParameters{
		Start:                       -1,
		End:                         -1,
		Sort:                        "id",
		Order:                       "ASC",
		Q:                           "",
		ID:                          []string{},
		IncludeAllIntoParseResource: false,
	}
	keys, ok := r.URL.Query()["id"]
	if ok && len(keys) >= 1 {
		params.ID = keys
	}
	keys, ok = r.URL.Query()["allinparse"]
	if ok && len(keys[0]) >= 1 {
		params.IncludeAllIntoParseResource = true
	}
	keys, ok = r.URL.Query()["_start"]
	if ok && len(keys[0]) >= 1 {
		val, err := strconv.Atoi(keys[0])
		if err == nil {
			params.Start = int64(val)
		}
	}
	keys, ok = r.URL.Query()["_end"]
	if ok && len(keys[0]) >= 1 {
		val, err := strconv.Atoi(keys[0])
		if err == nil {
			params.End = int64(val)
		}
	}
	if params.Start > params.End {
		params.Start, params.End = params.End, params.Start
	}
	keys, ok = r.URL.Query()["_sort"]
	if ok && len(keys[0]) >= 1 {
		params.Sort = keys[0]
	}
	keys, ok = r.URL.Query()["_order"]
	if ok && len(keys[0]) >= 1 {
		params.Order = keys[0]
	}
	keys, ok = r.URL.Query()["q"]
	if ok && len(keys[0]) >= 1 {
		params.Q = keys[0]
	}
	return params
}
