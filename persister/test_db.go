package persister

import (
	"context"
	"database/sql"
)

var _ DBTX = (*TestDB)(nil)

type TestDB struct {
}

func (db TestDB) ExecContext(context.Context, string, ...interface{}) (sql.Result, error) {
	return nil, nil
}
func (db TestDB) PrepareContext(context.Context, string) (*sql.Stmt, error) {
	return nil, nil
}
func (db TestDB) QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error) {
	return nil, nil
}
func (db TestDB) QueryRowContext(context.Context, string, ...interface{}) *sql.Row {
	return nil
}
