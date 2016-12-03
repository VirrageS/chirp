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
	CreatedAt     time.Time
	LastLogin     time.Time
	Active        bool
	Name          string
	AvatarUrl     sql.NullString
	Following     bool
}

type PublicUser struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Name      string `json:"name"`
	AvatarUrl string `json:"avatar_url"`
	Following bool   `json:"following"`
}