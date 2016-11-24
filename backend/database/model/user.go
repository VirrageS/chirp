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
}

type PublicUser struct {
	ID        int64
	Username  string
	Name      string
	AvatarUrl sql.NullString
}
