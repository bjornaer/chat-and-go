package main

import (
	"github.com/jinzhu/gorm"
)

// Message defined what is written and from which user
type Message struct {
	//gorm.Model
	ID        int    `json:"id" gorm:"primary_key"`
	Timestamp string `json:"timestamp"` // fix to work with real timestamp
	Username  string `json:"username"`  //gorm:"foreign_key"
	Content   string `json:"content"`
	Email     string `json:"email"`
}

type Messages []Message

type User struct {
	//gorm.Model
	ID       int    `gorm:"primary_key"`
	Username string `json:"username" db:"username"`
	Email    string `json:"email"`
}

type Server struct {
	db        *gorm.DB
	broadcast chan Message // broadcast channel
}
