package persister

import (
	_ "embed"
)

//go:generate sqlc generate -f ./sqlc.yaml
//go:generate ifacemaker -f *.go -s Queries -i DataStoreQueries -p persister -o query.iface.go
//go:generate ifacemaker -f *.go -s Queries -i DataStoreQueries -p app -o ../app/query.iface.go

//go:embed schema.sql
var ddl string

func DDL() string {
	return ddl
}
