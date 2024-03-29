package persister

import (
	_ "embed"
)

// IMPORTANT: Required version of sqlc incorporates this PR: https://github.com/sqlc-dev/sqlc/pull/3135
//go:generate sqlc generate -f ./sqlc.yaml
//go:generate ifacemaker -f *.go -s Queries -i DataStoreQueries -p persister -o query.iface.go
//go:generate moq -rm -skip-ensure -with-resets -pkg app -out ../app/query.stub.go . DataStoreQueries:DataStoreQueriesStub

//go:embed schema.sql
var ddl string

func DDL() string {
	return ddl
}
