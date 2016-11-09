package model

import "time"

type Tweet struct {
	ID        int64     `json:"id"`
	Author    User      `json:"author"`
	Likes     int64     `json:"likes"`
	Retweets  int64     `json:"retweets"`
	CreatedAt time.Time `json:"created_at"`
	Content   string    `json:"content"`
	Liked     bool      `json:"liked"`
	Retweeted bool      `json:"retweeted"`
}

type NewTweetContent struct {
	Content string `json:"content" binding:"required"`
}

type NewTweet struct {
	AuthorID int64
	Content  string
}
