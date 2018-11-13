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

// Configure the upgrader --> this will let us upgrade the http to a webSocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	var srvr Server
	srvr.broadcast = make(chan Message)
	db, err := gorm.Open("mysql", "root:testpass@tcp(db:3306)/gochat")
	//db, err := gorm.Open("mysql", "root:testpass@tcp(127.0.0.1:8989)/challenge")
	defer db.Close()
	if err != nil {
		log.Fatal("unable to connect to DB", err)
	}
	srvr.db = db
	srvr.db.AutoMigrate(&User{}, &Message{}) // move db stuff to a file

	rtr := srvr.SetupRoutes()
	/*
		Start listening for incoming chat messages from the broadcast channel
		and pass them to clients over their respective WebSocket connection.
	*/
	go srvr.HandleMessages()

	// Start the server on localhost port 8000 and log any errors
	log.Println("http server started on :8000")
	err = http.ListenAndServe(":8000", rtr)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
