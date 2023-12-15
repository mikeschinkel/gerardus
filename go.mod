module github.com/mikeschinkel/gerardus

go 1.21

replace github.com/mikeschinkel/go-typegen => ../go-typegen

replace github.com/mikeschinkel/go-diffator => ../go-diffator

replace github.com/mikeschinkel/go-serr => ../go-serr

replace github.com/mikeschinkel/go-lib => ../go-lib

require (
	github.com/google/go-cmp v0.6.0
	github.com/mattn/go-sqlite3 v1.14.18
	github.com/mikeschinkel/go-lib v0.0.0-00010101000000-000000000000
	github.com/mikeschinkel/go-serr v0.0.0-00010101000000-000000000000
	github.com/mikeschinkel/go-typegen v0.0.0-00010101000000-000000000000
	golang.org/x/mod v0.14.0
	golang.org/x/sync v0.5.0
)

require github.com/mikeschinkel/go-diffator v0.0.0-00010101000000-000000000000 // indirect
