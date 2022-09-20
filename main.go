package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"golang-api-server/config"
	"golang-api-server/handler"
	"golang-api-server/util"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"runtime"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/acme/autocert"
)

func Initialize() {
	// from Executable Directory
	ex, _ := os.Executable()
	fmt.Println("[DEBUG] Executable Dir:", filepath.Dir(ex))

	// Current working directory
	dir, _ := os.Getwd()
	fmt.Println("[DEBUG] Current Working Dir:", dir)

	// Relative on runtime DIR:
	_, b, _, _ := runtime.Caller(0)
	d1 := path.Join(path.Dir(b))
	fmt.Println("[DEBUG] Relative Dir", d1)
}

func main() {
	Initialize()

	l := log.New(os.Stdout, "golang-api: ", log.LstdFlags)
	v := util.NewValidation()

	// create a new serve mux and register the handlers
	sm := mux.NewRouter()
	// Specify the API version, in this case API version 1
	v1 := sm.PathPrefix("/v1").Subrouter()

	// create the handlers
	productHandler := handler.NewProductHandler(l, v)

	// GET Product
	getProductRouter := v1.Methods(http.MethodGet).Subrouter()
	getProductRouter.HandleFunc("/product", productHandler.GetProduct)
	// POST Product
	postProductRouter := v1.Methods(http.MethodPost).Subrouter()
	postProductRouter.HandleFunc("/product", productHandler.CreateProduct)
	// PUT Product
	putProductRouter := v1.Methods(http.MethodPut).Subrouter()
	putProductRouter.HandleFunc("/product", productHandler.UpdateProduct)
	// DELETE Product
	deleteProductRouter := v1.Methods(http.MethodDelete).Subrouter()
	deleteProductRouter.HandleFunc("/product", productHandler.DeleteProduct)

	BeginServer(sm, l)
}

func BeginServer(sm *mux.Router, l *log.Logger) {
	/* BEGIN::DEVELOPMENT MODE */
	if config.ENV == "dev" {
		devServer := &http.Server{
			Addr:         ":9090",           // configure the bind address to 9090 (for development)
			Handler:      sm,                // set the default handler
			ErrorLog:     l,                 // set the logger for the server
			ReadTimeout:  5 * time.Second,   // max time to read request from the client
			WriteTimeout: 10 * time.Second,  // max time to write response to the client
			IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
		}

		go func() {
			l.Println("[DEV] Server is starting on port 9090")
			err := devServer.ListenAndServe()
			if err != nil {
				l.Fatal(err)
			}
		}()

		// trap sigterm or interupt and gracefully shutdown the server
		c := make(chan os.Signal)
		signal.Notify(c, os.Interrupt)
		signal.Notify(c, os.Kill)

		// Block until a signal is received.
		sig := <-c
		l.Println("Received terminate, graceful shutdown", sig)

		// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
		ctx, ctxCancel := context.WithTimeout(context.Background(), 30*time.Second)
		devServer.Shutdown(ctx)
		defer ctxCancel()
	}
	/* END::DEVELOPMENT MODE */

	/* BEGIN::PRODUCTION MODE */
	if config.ENV == "prod" {
		// Initialize certbot as certification manager
		certManager := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist("blablabal.com"), // ! edit this
			Cache:      autocert.DirCache("certs"),              //Folder for storing certificates
		}

		// create a new server for production
		prodServer := &http.Server{
			Addr:    ":443", // configure the bind address to 443 (HTTPS)
			Handler: sm,     // set the default handler
			TLSConfig: &tls.Config{
				GetCertificate: certManager.GetCertificate,
			},
			ErrorLog:     l,                 // set the logger for the server
			ReadTimeout:  5 * time.Second,   // max time to read request from the client
			WriteTimeout: 10 * time.Second,  // max time to write response to the client
			IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
		}

		go http.ListenAndServe(":80", certManager.HTTPHandler(nil))
		prodServer.ListenAndServeTLS("", "")
	}
	/* END::PRODUCTION MODE */
}
