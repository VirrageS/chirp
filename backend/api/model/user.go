package model

import "time"

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	LastLogin time.Time `json:"last_login"`
	Name      string    `json:"name"`
	AvatarUrl string    `json:"avatar_url"`
	Following bool      `json:"following"`
}
