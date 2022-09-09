package handler

import (
	"golang-api-server/util"
	"net/http"
)

func ProductHandler(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusNoContent)
}

func GetProduct(writer http.ResponseWriter, request *http.Request) {
	// ToJSON(&Product{}, writer)
	util.ToJSON(map[string]interface{}{
		"msg": "product list",
	}, writer)
}
