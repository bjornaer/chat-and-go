package main

import (
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/websocket"
	"github.com/jinzhu/gorm"
)

/* The value here in clients isn't actually needed but we are using a map
because it is easier than an array to append and delete items. */
var clients = make(map[*websocket.Conn]bool) // connected clients
//var broadcast = make(chan Message)           // broadcast channel

// Configure the upgrader --> this will let us upgrade the http to a webSocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	var srvr Server
	var err error
	srvr.broadcast = make(chan Message)
	srvr.db, err = gorm.Open("mysql", "root:testpass@tcp(db:3306)/challenge")
	defer srvr.db.Close()
	if err != nil {
		log.Fatal("unable to connect to DB", err)
	}
	srvr.db.AutoMigrate(&User{}, &Message{}) // move db stuff to a file
	/*
		if !srvr.db.HasTable(&User{}) {
			srvr.db.CreateTable(&User{})
		} else if !srvr.db.HasTable(&Message{}) {
			srvr.db.CreateTable(&Message{})
		}
	*/

	srvr.SetupRoutes()
	/*
		Start listening for incoming chat messages from the broadcast channel
		and pass them to clients over their respective WebSocket connection.
	*/
	go srvr.HandleMessages()

	// Start the server on localhost port 8000 and log any errors
	log.Println("http server started on :8000")
	err = http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
