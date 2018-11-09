package main

import (
	"github.com/jinzhu/gorm"
)

// Message defined what is written and from which user
type Message struct {
	gorm.Model
	Username string `json:"username";gorm:"foreign_key"`
	Message  string `json:"message"`
}

type Messages []Message

type Chats struct {
	Mssgs Messages `json:messages`
}

type User struct {
	gorm.Model
	Username string `gorm: "unique"`
}

type Server struct {
	db        *gorm.DB
	broadcast chan Message // broadcast channel
}
