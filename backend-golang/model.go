package main

import (
	"github.com/jinzhu/gorm"
)

// Message defined what is written and from which user
type Message struct {
	//gorm.Model
	ID int `json:"id" gorm:"primary_key"`
	//Timestamp time.Time `json:"timestamp"`
	Username string `json:"username"` //gorm:"foreign_key"
	Content  string `json:"content"`
}

type Messages []Message

type User struct {
	//gorm.Model
	ID       int    `gorm:"primary_key"`
	Username string `json:"username" db:"username" gorm:"unique"`
}

type Server struct {
	db        *gorm.DB
	broadcast chan Message // broadcast channel
}
