package model

import (
	"database/sql"
)

type PublicUser struct {
	ID        int64
	Username  string
	Name      string
	AvatarUrl sql.NullString
	Following bool
}
