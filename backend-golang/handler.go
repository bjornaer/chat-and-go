package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func (s Server) TestConnection(w http.ResponseWriter, r *http.Request) {
	result := "testing"

	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{
		"result":  result,
		"backend": "go",
	}); err != nil {
		log.Panic(err)
	}
}

func (s Server) FetchHistory(w http.ResponseWriter, r *http.Request) {
	/*
		var result string
		if err := s.db.QueryRow(`SELECT col FROM test`).Scan(&result); err != nil {
			log.Panic(err)
		}

		w.Header().Add("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]string{
			"result":  result,
			"backend": "go",
		}); err != nil {
			log.Panic(err)
		}
	*/
}

func (s Server) Login(w http.ResponseWriter, r *http.Request) {
	// read request body to get user json
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// unmarshal user json into usr variable
	var usr User
	err = json.Unmarshal(b, &usr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// find if user exists or create record if not, assign it to loggedUsr
	var loggedUsr User
	s.db.Where("username = ?", usr.Username).FirstOrCreate(&loggedUsr)
	// respond with usename and id
	output, err := json.Marshal(loggedUsr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(output)
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
		// {code for that here}
		// Send it out to every client that is currently connected
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
