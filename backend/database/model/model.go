package model

import "time"

type Tweet struct {
	ID           int64
	AuthorID     int64
	LikeCount    int64
	RetweetCount int64
	CreatedAt    time.Time
	Content      string
}

type User struct {
	ID        int64
	Name      string
	Username  string
	Email     string
	CreatedAt time.Time
}
