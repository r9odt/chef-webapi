package handler

import (
	"fmt"
	"net/http"

	"github.com/JIexa24/chef-webapi/httpserver/controller"
)

func emptyJSON() string {
	return "{}"
}

func setJSONHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

func setXTotalCountHeader(w http.ResponseWriter, length int64) {
	w.Header().Set("X-Total-Count",
		fmt.Sprintf("%d", length))
}

func getStartAndEndFromParams(params *controller.APIParameters,
	length int64) (int64, int64) {
	start := int64(0)
	end := int64(length)
	if params.Start >= 0 {
		start = params.Start
	}
	if start < 0 {
		start = 0
	}
	if params.End >= 0 {
		end = params.End
	}
	if end > length {
		end = length
	}
	return start, end
}
