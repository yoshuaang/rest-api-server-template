package main

import (
	"golang-api-server/handler"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	v1 := router.PathPrefix("/v1").Subrouter()

	getProductRouter := v1.Methods(http.MethodGet).Subrouter()
	getProductRouter.HandleFunc("/product", handler.GetProduct).Methods("GET")
	log.Fatal(http.ListenAndServe(":9090", router))
}
