package model

import "time"

type Tweet struct {
	ID        int64
	AuthorID  int64
	Likes     int64
	Retweets  int64
	CreatedAt time.Time
	Content   string
}
