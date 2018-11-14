package main

import (
	"time"

	"github.com/jinzhu/gorm"
)

// my tables/structs are not normalized. Should normalize them and use smarter queries. But for now, this works.

// Message defined what is written and from which user
type Message struct {
	//gorm.Model
	ID        int       `json:"id" gorm:"primary_key"`
	Timestamp time.Time `json:"timestamp"`
	Username  string    `json:"username"`
	Content   string    `json:"content"`
	Email     string    `json:"email"`
}

// Messages array of Message
type Messages []Message

// User is a person
type User struct {
	//gorm.Model
	ID              int       `json:"id" gorm:"primary_key"`
	Username        string    `json:"username" db:"username"`
	Email           string    `json:"email" gorm:"index"`
	LastInteraction time.Time `json:"lastInteraction"`
}

// Handler holds the context so everyone uses the same channel to talk to the db and is aware of the broadcast
type Handler struct {
	db        *gorm.DB
	broadcast chan Message // broadcast channel
}
