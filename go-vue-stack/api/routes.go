package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

// where frontend stuff awaits
const frontEntry = "../public"

// SetupRoutes says who handles what and where
func (h Handler) SetupRoutes() *mux.Router {
	r := mux.NewRouter()
	// create the login route based on the api-attempt!
	r.HandleFunc("/login", h.Login).Methods("POST")
	r.HandleFunc("/history", h.FetchHistory).Queries("oldest", "{oldest}").Methods("GET")
	r.HandleFunc("/newMessages", h.FetchNewMessages).Queries("id", "{id}").Methods("GET")
	r.HandleFunc("/test", h.TestConnection).Methods("GET")
	// Configure websocket route
	r.HandleFunc("/ws", h.HandleConnections)
	// Create a simple file server
	fs := http.FileServer(http.Dir(frontEntry))
	r.PathPrefix("/").Handler(fs)
	return r
}
