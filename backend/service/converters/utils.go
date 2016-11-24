package converters

import "database/sql"

func toSqlNullString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: s != ""}
}
