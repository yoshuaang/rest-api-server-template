package main

import (
	"golang-api-server/handler"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter()
	v1 := router.PathPrefix("/v1").Subrouter()

	getProductRouter := v1.Methods(http.MethodGet).Subrouter()
	getProductRouter.HandleFunc("/product", handler.GetProduct).Methods("GET")

	postProductRouter := v1.Methods(http.MethodPost).Subrouter()
	postProductRouter.HandleFunc("/product", handler.CreateProduct).Methods("POST")

	return router
}

func main() {
	router := Router()
	log.Fatal(http.ListenAndServe(":9090", router))
}
