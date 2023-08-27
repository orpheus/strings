package sqldb

import "database/sql"

// NewNullString converts empty strings to NullString so
// that the DB receives NULL instead of an empty string
// if string is empty.
// https://stackoverflow.com/questions/40266633/golang-insert-null-into-sql-instead-of-empty-string
func NewNullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}
