package app

import (
	"context"
	"database/sql"

	"github.com/mikeschinkel/gerardus/persister"
)

type DataStore interface {
	Open() error
	Query(ctx context.Context, sql string) error
	DB() *sql.DB
	Initialize(ctx context.Context) error
	Queries() persister.DataStoreQueries
}
