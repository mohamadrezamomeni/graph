package sqlite

import (
	"errors"

	"github.com/mattn/go-sqlite3"
)

func IsDuplicateError(err error) bool {
	var sqliteErr sqlite3.Error

	if !errors.As(err, &sqliteErr) {
		return false
	}

	if sqliteErr.Code != sqlite3.ErrConstraint {
		return false
	}

	if sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique ||
		sqliteErr.ExtendedCode == sqlite3.ErrConstraintPrimaryKey {
		return true
	}

	return false
}
