package db

import (
	"errors"

	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

func IsUniqueViolation(err error) bool {
	var sqliteError *sqlite.Error
	if errors.As(err, &sqliteError) {
		return sqliteError.Code() == sqlite3.SQLITE_CONSTRAINT || sqliteError.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE
	}

	return false
}
