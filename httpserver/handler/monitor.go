package handler

import (
	"fmt"
	"net/http"
)

func monitorReadyAPIHandler(writer http.ResponseWriter, request *http.Request) {
	setJSONHeader(writer)
	fmt.Print(writer, emptyJSON())
}

func monitorHealthAPIHandler(writer http.ResponseWriter, request *http.Request) {
	setJSONHeader(writer)
	fmt.Print(writer, emptyJSON())
}
