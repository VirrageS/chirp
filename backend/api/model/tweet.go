package model

import "time"

type Tweet struct {
	ID           int64     `json: "author_id"`
	Author       User      `json: "user"`
	LikeCount    int64     `json: "likes"`
	RetweetCount int64     `json: "retweets"`
	CreatedAt    time.Time `json: "created_at"`
	Content      string    `json: "content"`
}

type NewTweet struct {
	AuthorID int64  `json: "id"`
	Content  string `json: "content"`
}
