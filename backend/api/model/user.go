package model

import "time"

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	LastLogin time.Time `json:"last_login"`
	Active    bool      `json:"active"`
	Name      string    `json:"name"`
	AvatarUrl string    `json:"avatar_url"`
	Following bool      `json:"following"`
}
