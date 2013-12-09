package main

import (
	"github.com/gorilla/mux"
	"github.com/r4mp/c3an/api"
	"log"
	"net/http"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU()) // use all CPU cores

	go receiveMessages()

	log.Println("Starting server...")

	m := mux.NewRouter()

	m.HandleFunc("/device/register", api.RegisterDevice).Methods("POST")
	m.HandleFunc("/device/unregister", api.UnregisterDevice).Methods("POST")

	m.HandleFunc("/notification/send", api.SendNotification).Methods("POST")

	// Everything else fails.
	m.HandleFunc("/{path:.*}", http.NotFound)

	log.Println("Now listening on port 8080")

	http.Handle("/", m)
	http.ListenAndServe(":8080", nil)
}
