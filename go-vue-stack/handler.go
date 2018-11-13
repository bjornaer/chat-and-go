package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jinzhu/gorm"
)

// TestConnection just asks for response. Nothing special!
func (s Server) TestConnection(w http.ResponseWriter, r *http.Request) {
	//result := "testing"

	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{
		"result":  "result",
		"backend": "go",
	}); err != nil {
		log.Panic(err)
	}
}

// FetchHistory asks for latest 100 messages as to start chat with a short history
func (s Server) FetchHistory(w http.ResponseWriter, r *http.Request) {
	var msgs []Message
	err := s.db.Order("ID desc").Limit(100).Find(&msgs).Error
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(map[string]Messages{
		"messages": msgs,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	//w.Write(output)
}

// Login is the only endpoint for registry and/or login
func (s Server) Login(w http.ResponseWriter, r *http.Request) {
	// read request body to get user json
	var usr User
	err := json.NewDecoder(r.Body).Decode(&usr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// find if user exists or create record if not, assign it to loggedUsr
	var loggedUsr User
	if err = s.db.Where("email = ?", usr.Email).First(&loggedUsr).Error; gorm.IsRecordNotFoundError(err) {
		s.db.Create(&usr)
		loggedUsr = usr
	}
	// respond with usename and id
	if usr.Username != loggedUsr.Username {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	err = json.NewEncoder(w).Encode(loggedUsr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	//w.Write(output)
}

// HandleConnections manages new socket connections and incoming messages
func (s Server) HandleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	// Make sure we close the connection when the function returns
	defer ws.Close()

	// Register our new client
	clients[ws] = true

	for {
		var msg Message
		/*
			Read in a new message as JSON and map it to a Message object.
			If there is some kind of error with reading from the socket,
			it's safe to assume the client has disconnected for some reason or another.
			We log the error and remove that client from our global "clients" map
			so we don't try to read from or send new messages to that client
		*/
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}
		// Send the newly received message to the broadcast channel
		s.broadcast <- msg
	}
}

// HandleMessages as the name suggests handles that. Sends boradcasted messages to everyone else!
func (s Server) HandleMessages() {
	for {
		// Grab the next message from the broadcast channel
		msg := <-s.broadcast
		// Store new message in the DB
		s.db.Create(&msg)
		// Send it out to every client that is currently connected
		for client := range clients {
			err := client.WriteJSON(&msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
