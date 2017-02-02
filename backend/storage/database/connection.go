package database

import (
	"database/sql"
)

// Connection is wrapper for sql.DB struct.
// Name is easier to reason about. It provides a lot clarity in the code.
type Connection struct {
	*sql.DB
}
