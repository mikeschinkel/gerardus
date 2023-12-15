package persister

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"log/slog"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mikeschinkel/gerardus/paths"
)

var dataStore *DataStore

type DataStore struct {
	*Queries
	Filepath string
	db       *sql.DB
}

func NewDataStore(dbFile string) *DataStore {
	return &DataStore{
		Filepath: dbFile,
	}
}

func (db *DataStore) Open() (err error) {
	db.db, err = sql.Open("sqlite3", db.Filepath)
	if err != nil {
		goto end
	}
	db.Queries = New(db.db)
end:
	return err
}

func (db *DataStore) Query(ctx context.Context, sql string) (err error) {
	_, err = db.db.ExecContext(ctx, sql)
	return err
}

func (db *DataStore) InitializeDataStore(ctx context.Context) (err error) {
	slog.Info("Initializing data store")
	err = db.Open()
	if err != nil {
		goto end
	}
	err = db.Query(ctx, DDL())
	if err != nil {
		goto end
	}
end:
	return err
}

func getDataStore(fp string) (ds *DataStore, err error) {
	absFP, err := paths.Absolute(fp)
	if err != nil {
		err = fmt.Errorf("error attempted to convert '%s' to an absolute path: %s; %w",
			fp, err)
		goto end
	}
	ds = NewDataStore(absFP)
end:
	return ds, err
}

func GetDataStore() *DataStore {
	if dataStore == nil {
		panic("DataStore not yet initialized")
	}
	return dataStore
}
