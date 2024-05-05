package models

import "time"

type User struct {
	Id        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Password  string    `json:"password,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
}
