package persister

import (
	_ "embed"
)

//go:generate sqlc generate -f ./sqlc.yaml
//go:generate ifacemaker -f query.sql.go -s Queries -i DataStoreQueries -p persister -o query.iface.go

//go:embed schema.sql
var ddl string

func DDL() string {
	return ddl
}
