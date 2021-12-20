package models

import (
	"time"
)

type User struct {
	ID        int64     `json:"id"`
	Firstname string    `json:"firstname"`
	Lastname  string    `json:"lastname"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Deleted   bool      `json:"deleted"`
	Nickname  string    `json:"nickname"`
	CreatedAt time.Time `json:"created_at"`
}
