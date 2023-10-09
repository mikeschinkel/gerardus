package main

import (
	"gerardus/cli"
	"gerardus/options"
	"gerardus/persister"
)

func init() {
	cli.RootCmd.AddFlag(&cli.Flag{
		Name:             "data",
		VarName:          "data_file",
		Usage:            "Data file (sqlite3)",
		Default:          persister.SqliteDB,
		SetStringValFunc: options.SetDataFile,
	})
}
