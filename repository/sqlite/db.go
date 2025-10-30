package sqlite

import (
	"database/sql"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	appError "github.com/mohamadrezamomeni/graph/pkg/error"
	"github.com/mohamadrezamomeni/graph/pkg/utils"
)

type SqliteDB struct {
	db *sql.DB
}

func (s *SqliteDB) Conn() *sql.DB {
	return s.db
}

func New(cfg *DBConfig) *SqliteDB {
	const scope = "sqlite.New"

	root, err := utils.GetRootOfProject()
	if err != nil {
		panic(
			appError.Wrap(err).
				Scope(scope).
				DeactiveWrite().
				Input(cfg).
				Errorf("failed to locate project root: %s", err.Error()),
		)
	}

	dbPath := filepath.Join(root, cfg.Path)
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		panic(
			appError.Wrap(err).
				Scope(scope).
				DeactiveWrite().
				Input(cfg).
				Errorf("failed to open SQLite database at path %s: %s", dbPath, err.Error()),
		)
	}
	db.Exec("PRAGMA foreign_keys = ON")

	return &SqliteDB{
		db: db,
	}
}
