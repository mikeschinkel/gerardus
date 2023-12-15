module github.com/mikeschinkel/gerardus/cmd

go 1.21

replace github.com/mikeschinkel/go-typegen => ../../go-typegen

replace github.com/mikeschinkel/go-diffator => ../../go-diffator

replace github.com/mikeschinkel/go-serr => ../../go-serr

replace github.com/mikeschinkel/go-lib => ../../go-lib

replace github.com/mikeschinkel/gerardus => ../../gerardus

require (
	github.com/mikeschinkel/gerardus v0.0.0-00010101000000-000000000000
	github.com/mikeschinkel/go-serr v0.0.0-20231130133231-53784ffdd4bf
)

require (
	github.com/mattn/go-sqlite3 v1.14.18 // indirect
	github.com/mikeschinkel/go-diffator v0.0.0-00010101000000-000000000000 // indirect
	github.com/mikeschinkel/go-lib v0.0.0-20231127221959-1bbf62bd0845 // indirect
	github.com/mikeschinkel/go-typegen v0.0.0-00010101000000-000000000000 // indirect
	golang.org/x/mod v0.14.0 // indirect
	golang.org/x/sync v0.5.0 // indirect
)
