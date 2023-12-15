package app

import (
	"reflect"

	"github.com/mikeschinkel/gerardus/cli"
	"github.com/mikeschinkel/gerardus/options"
	"github.com/mikeschinkel/gerardus/persister"
)

func init() {
	cli.RootCmd.
		AddFlag(&cli.Flag{
			Switch: "data",
			Arg: &cli.Arg{
				Name:         "data_file",
				Usage:        "Data file (sqlite3)",
				Type:         reflect.String,
				Default:      persister.SqliteDB,
				SetValueFunc: options.SetDataFile,
			}})
}
