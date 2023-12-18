package persister

import (
	"context"
	"database/sql"
	_ "embed"
	"log/slog"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mikeschinkel/gerardus/paths"
)

var _ DataStore = (*SqliteDataStore)(nil)

type SqliteDataStore struct {
	queries  DataStoreQueries
	Filepath string
	db       *sql.DB
}

func NewSqliteDataStore(dbFile string) DataStore {
	return &SqliteDataStore{
		Filepath: dbFile,
	}
}

func (ds *SqliteDataStore) Initialize(ctx context.Context) (err error) {
	slog.Info("Initializing data store")

	absFP, err := paths.Absolute(ds.Filepath)
	if err != nil {
		err = ErrFailedConvertToAbsPath.Err(err, "filepath", ds.Filepath)
		goto end
	}
	ds.Filepath = absFP

	err = ds.Open()
	if err != nil {
		goto end
	}
	err = ds.Query(ctx, DDL())
	if err != nil {
		goto end
	}
end:
	return err
}

func (db *SqliteDataStore) Open() (err error) {
	db.db, err = sql.Open("sqlite3", db.Filepath)
	if err != nil {
		goto end
	}
	db.queries = New(db.db)
end:
	return err
}

func (db *SqliteDataStore) Query(ctx context.Context, sql string) (err error) {
	_, err = db.db.ExecContext(ctx, sql)
	return err
}

func (db *SqliteDataStore) Queries() DataStoreQueries {
	return db.queries
}

func (db *SqliteDataStore) DB() *sql.DB {
	return db.db
}
