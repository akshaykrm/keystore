package db

import (
	"database/sql"
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

func IsNotFound(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}
