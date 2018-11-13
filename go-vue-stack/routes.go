package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

// getenv variable to use react frontend or vue
const frontEntry = "./public"

// SetupRoutes says who handles what and where
func (s Server) SetupRoutes() *mux.Router {
	r := mux.NewRouter()
	// create the login route based on the api-attempt!
	r.HandleFunc("/login", s.Login).Methods("POST")
	r.HandleFunc("/history", s.FetchHistory).Methods("GET")
	r.HandleFunc("/test", s.TestConnection).Methods("GET")
	// Configure websocket route
	r.HandleFunc("/ws", s.HandleConnections)
	// Create a simple file server
	fs := http.FileServer(http.Dir(frontEntry))
	r.PathPrefix("/").Handler(fs)
	return r
}
