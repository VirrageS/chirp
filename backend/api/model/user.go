package model

import "time"

type User struct {
	ID        int64     `json: "id"`
	Name      string    `json: "name"`
	Username  string    `json: "username"`
	Email     string    `json: "email"`
	CreatedAt time.Time `json: "created_at"`
}

type NewUser struct {
	Name     string `json: "name"`
	Username string `json: "username"`
	Email    string `json: "email"`
}
