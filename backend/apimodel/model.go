package apimodel

import (
	"time"
)

type Tweet struct {
	Id           	int64 `json: "author_id"`
	Author       	User `json: "user"`
	LikeCount    	int64 `json: "likes"`
	RetweetCount 	int64 `json: "retweets"`
	CreatedAt    	time.Time `json: "created_at"`
	Content      	string `json: "content"`
}

type NewTweet struct {
	AuthorId	int64 `json: "id"`
	Content		string `json: "content"`
}

type User struct {
	Id 		int64 `json: "id"`
	Name 		string `json: "name"`
	Username 	string `json: "username"`
	Email	 	string `json: "email"`
	CreatedAt 	time.Time `json: "email"`
}