package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

const frontEntry = "../frontend-react/public"

func (s Server) SetupRoutes() {
	r := mux.NewRouter()
	// create the login route based on the api-attempt!
	r.HandleFunc("/login", s.Login).Methods("POST")
	r.HandleFunc("/history", s.FetchHistory).Methods("GET") // implement
	r.HandleFunc("/test", s.TestConnection).Methods("GET")
	// Create a simple file server
	fs := http.FileServer(http.Dir(frontEntry))
	r.Handle("/", fs).Methods("GET")
	// Configure websocket route
	r.HandleFunc("/ws", s.HandleConnections)
}
