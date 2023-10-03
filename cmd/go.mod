module gerardus/cmd

go 1.21

replace gerardus => ../.

require (
	golang.org/x/sync v0.3.0
	gerardus v0.0.0-00010101000000-000000000000
)

require github.com/mattn/go-sqlite3 v1.14.17 // indirect
