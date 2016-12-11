package model

import "time"

type Tweet struct {
	ID           int64       `json:"id"`
	Author       *PublicUser `json:"author"`
	LikeCount    int64       `json:"like_count"`
	RetweetCount int64       `json:"retweet_count"`
	CreatedAt    time.Time   `json:"created_at"`
	Content      string      `json:"content"`
	Liked        bool        `json:"liked"`
	Retweeted    bool        `json:"retweeted"`
}

type NewTweetContent struct {
	Content string `json:"content" binding:"required"`
}

type NewTweet struct {
	AuthorID int64  `json:"-"`
	Content  string `json:"content" binding:"required"`
}
