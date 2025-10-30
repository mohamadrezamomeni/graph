package contact

import (
	"github.com/mohamadrezamomeni/graph/repository/sqlite"
)

type Contact struct {
	db *sqlite.SqliteDB
}

func New(conn *sqlite.SqliteDB) *Contact {
	return &Contact{
		db: conn,
	}
}
