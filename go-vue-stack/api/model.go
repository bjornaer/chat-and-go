package main

import (
	"time"

	"github.com/jinzhu/gorm"
)

// Message defined what is written and from which user, adding the field User and UserID allows to set a 'belongs to' type of relation
type Message struct {
	//gorm.Model
	ID        int       `json:"id" gorm:"primary_key"`
	Timestamp time.Time `json:"timestamp"`
	Content   string    `json:"content"`
	User      User      `json:"user"`
	UserID    int       `json:"userId"`
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
