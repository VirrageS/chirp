package model

import (
	"database/sql"
	"time"
)

type User struct {
	ID            int64
	TwitterToken  sql.NullString
	FacebookToken sql.NullString
	GoogleToken   sql.NullString
	Username      string
	Password      string
	Email         string
	CreatedAt     *time.Time
	LastLogin     *time.Time
	Active        bool
	Name          string
	AvatarUrl     sql.NullString
	FollowerCount int64
	FolloweeCount int64
	Following     bool
}

type PublicUser struct {
	ID            int64  `json:"id"`
	Username      string `json:"username"`
	Name          string `json:"name"`
	AvatarUrl     string `json:"avatar_url"`
	FollowerCount int64  `json:"follower_count"`
	FolloweeCount int64  `json:"followee_count"`
	Following     bool   `json:"following"`
}

type UserGoogle struct {
    Sub string `json:"sub"`
    Name string `json:"name"`
    GivenName string `json:"given_name"`
    FamilyName string `json:"family_name"`
    Picture string `json:"picture"`
    Email string `json:"email"`
    EmailVerified string `json:"email_verified"`
}
