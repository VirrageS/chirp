package model

import "time"

// TODO: Maybe it should contain whole user struct or there should be another model struct that has whole User object embedded
type Tweet struct {
	ID        int64
	AuthorID  int64
	Likes     int64
	Retweets  int64
	CreatedAt time.Time
	Content   string
}

type TweetWithAuthor struct {
	ID        int64
	Author    *PublicUser
	Likes     int64
	Retweets  int64
	CreatedAt time.Time
	Content   string
}
