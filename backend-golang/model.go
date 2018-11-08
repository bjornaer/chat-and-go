package main

import "database/sql"

// Message defined what is written and from which user
type Message struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Message  string `json:"message"`
}

type Server struct {
	db        *sql.DB
	broadcast chan Message // broadcast channel
}
