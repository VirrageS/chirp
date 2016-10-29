package model

import "time"

type User struct {
	ID        int64
	Name      string
	Username  string
	Email     string
	CreatedAt time.Time
}
