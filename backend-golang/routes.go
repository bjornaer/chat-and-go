package main

import "net/http"

const frontEntry = "../frontend-react/public"

func (s Server) SetupRoutes() {
	http.HandleFunc("/test", s.TestConnection)
	// Create a simple file server
	fs := http.FileServer(http.Dir(frontEntry))
	http.Handle("/", fs)
	// Configure websocket route
	http.HandleFunc("/ws", s.HandleConnections)
}
