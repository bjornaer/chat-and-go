package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func (s Server) TestConnection(w http.ResponseWriter, r *http.Request) {
	//result := "testing"

	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{
		"result":  "result",
		"backend": "go",
	}); err != nil {
		log.Panic(err)
	}
	/* testUsr := User{Username: "Chatengo"}
	s.db.Create(&testUsr)
	testMsg := Message{Username: testUsr.Username, Content: "Hello There!"}
	s.db.Create(&testMsg)
	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"testMsg": testMsg,
		"testUsr": testUsr,
	}); err != nil {
		log.Panic(err)
	} */
}

func (s Server) FetchHistory(w http.ResponseWriter, r *http.Request) {
	/*
		TODO get paginated result for all messages
		serve it for chat
	*/
	var msgs []Message
	err := s.db.Order("ID desc").Limit(100).Find(&msgs).Error
	/* err := json.NewEncoder(w).Encode(map[string]Messages{
		"messages": msgs,
	}) */
	//var msg Message
	//s.db.Last(&msg)
	err = json.NewEncoder(w).Encode(msgs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	//w.Write(output)
}

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
	s.db.Where("username = ?", usr.Username).FirstOrCreate(&loggedUsr)
	fmt.Println("USER SHOULD BE: " + loggedUsr.Username)
	// respond with usename and id
	err = json.NewEncoder(w).Encode(loggedUsr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	//w.Write(output)
}

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
