package migrate

import (
	"database/sql"
	"path/filepath"

	"github.com/mohamadrezamomeni/graph/repository/sqlite"

	momoError "github.com/mohamadrezamomeni/graph/pkg/error"
	"github.com/mohamadrezamomeni/graph/pkg/utils"

	momoLogger "github.com/mohamadrezamomeni/graph/pkg/log"

	migrate "github.com/rubenv/sql-migrate"
)

type Migrator struct {
	dialect    string
	path       string
	migrations *migrate.FileMigrationSource
}

func New(cfg *sqlite.DBConfig) *Migrator {
	root, _ := utils.GetRootOfProject()
	migrations := &migrate.FileMigrationSource{
		Dir: filepath.Join(root, "./repository/sqlite/migrations"),
	}

	return &Migrator{
		path:       filepath.Join(root, cfg.Path),
		dialect:    "sqlite3",
		migrations: migrations,
	}
}

func (m *Migrator) UP() {
	scope := "migration.up"
	db, err := sql.Open(m.dialect, m.path)
	if err != nil {
		panic(
			momoError.
				Wrap(err).
				Scope(scope).
				DeactiveWrite().
				Errorf("error to connect db"))
	}

	n, err := migrate.Exec(db, m.dialect, m.migrations, migrate.Up)
	if err != nil {
		panic(momoError.
			Wrap(err).
			DeactiveWrite().
			Scope(scope).
			Errorf("unable to apply migrations: %s", err.Error()))
	}
	momoLogger.Infof("Applied %d migrations!", n)
}

func (m *Migrator) DOWN() {
	scope := "migration.down"

	db, err := sql.Open(m.dialect, m.path)
	if err != nil {
		panic(
			momoError.
				Wrap(err).
				DeactiveWrite().
				Scope(scope).
				Errorf("error to connect db"),
		)
	}

	n, err := migrate.Exec(db, m.dialect, m.migrations, migrate.Down)
	if err != nil {
		panic(
			momoError.
				Wrap(err).
				DeactiveWrite().
				Scope(scope).
				Errorf("unable to undo migrations: %s", err.Error()))
	}
	momoLogger.Infof("undo %d migrations!", n)
}
