package main

import (
	"golang-api-server/handler"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/product", handler.ProductHandler).Methods("GET")
	log.Fatal(http.ListenAndServe(":9090", router))
}
