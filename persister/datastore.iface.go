package persister

import (
	"context"
	"database/sql"
)

type DataStore interface {
	Open() error
	Query(ctx context.Context, sql string) error
	DB() *sql.DB
	Initialize(ctx context.Context) error
	Queries() DataStoreQueries
	SetQueries(DataStoreQueries)
}
