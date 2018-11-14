package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
)

// TestConnection just asks for response. Nothing special!
func (h Handler) TestConnection(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{
		"result":  "result",
		"backend": "go",
	}); err != nil {
		log.Panic(err)
	}
}

// FetchHistory asks for latest 100 messages as to start chat with a short history, but on next requests it starts going back in time.
func (h Handler) FetchHistory(w http.ResponseWriter, r *http.Request) {
	var msgs []Message
	var err error
	oldestLoadedID := r.FormValue("oldest")
	messageAmount := r.FormValue("quantity")
	id, _ := strconv.ParseInt(oldestLoadedID, 10, 64)
	quantity, _ := strconv.ParseInt(messageAmount, 10, 64)
	if id < 0 {
		err = h.db.Order("ID desc").Limit(quantity).Find(&msgs).Error
	} else {
		err = h.db.Where("ID < ?", oldestLoadedID).Order("ID desc").Limit(quantity).Find(&msgs).Error
	}
	// Modify to return paginated result!

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	setHeaderOKStatus(w)
	err = json.NewEncoder(w).Encode(map[string]Messages{
		"messages": msgs,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	/* w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK) */
}

// FetchNewMessages asks for latest 10 messages since last registered interaction
func (h Handler) FetchNewMessages(w http.ResponseWriter, r *http.Request) {
	var msgs []Message
	var user User
	userID := r.FormValue("id")
	id, _ := strconv.ParseInt(userID, 10, 64)
	err := h.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = h.db.Where("timestamp > ?", user.LastInteraction).Order("timestamp asc").Limit(10).Find(&msgs).Error
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	setHeaderOKStatus(w)
	err = json.NewEncoder(w).Encode(map[string]Messages{
		"messages": msgs,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	/* w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK) */

}

// Login is the only endpoint for registry and/or login
func (h Handler) Login(w http.ResponseWriter, r *http.Request) {
	// read request body to get user json
	var usr User
	err := json.NewDecoder(r.Body).Decode(&usr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// find if user exists or create record if not, assign it to loggedUsr
	var loggedUsr User
	if err = h.db.Where("email = ?", usr.Email).First(&loggedUsr).Error; gorm.IsRecordNotFoundError(err) {
		// if the user did not exist, we set last interaction as now. If it did exist, let's keep the last value!
		usr.LastInteraction = time.Now()
		h.db.Create(&usr)
		loggedUsr = usr
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// respond with usename and id
	if usr.Username != loggedUsr.Username {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	setHeaderOKStatus(w)
	err = json.NewEncoder(w).Encode(loggedUsr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	/* w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK) */
}

// HandleConnections manages new socket connections and incoming messages
func (h Handler) HandleConnections(w http.ResponseWriter, r *http.Request) {
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
		var user User
		/*
			Read in a new message as JSON and map it to a Message object.
			If there is some kind of error with reading from the socket,
			it's safe to assume the client has disconnected for some reason or another.
			We log the error and remove that client from our global "clients" map
			so we don't try to read from or send new messages to that client
		*/
		err := ws.ReadJSON(&msg)
		msg.Timestamp = time.Now()
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}
		// Store new message in the DB
		err = h.db.Create(&msg).Error
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// update lastInteraction for user
		err = h.db.Where("email = ?", msg.Email).First(&user).Error
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = h.db.Model(&user).Update("LastInteraction", time.Now()).Error
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Send the newly received message to the broadcast channel
		h.broadcast <- msg
	}
}

// HandleMessages as the name suggests handles that. Sends boradcasted messages to everyone else!
func (h Handler) HandleMessages() {
	for {
		// Grab the next message from the broadcast channel
		msg := <-h.broadcast
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

func setHeaderOKStatus(w http.ResponseWriter) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
}
