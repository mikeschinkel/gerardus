package app

import (
	"reflect"

	"github.com/mikeschinkel/gerardus/cli"
	"github.com/mikeschinkel/gerardus/options"
	"github.com/mikeschinkel/gerardus/paths"
	"github.com/mikeschinkel/gerardus/persister"
)

func init() {
	cli.RootCmd.
		AddFlag(cli.Flag{
			Switch: "data",
			Arg: cli.Arg{
				Name:         "data_file",
				Usage:        "Data file (sqlite3)",
				Type:         reflect.String,
				Default:      persister.SqliteDB,
				ValidateFunc: Root.validateDataFile,
				SetValueFunc: options.SetDataFile,
			}})
}

func (a *App) validateDataFile(ctx Context, file any, arg *cli.Arg) (err error) {
	fileName := file.(string)
	exists, err := paths.FileExists(fileName)
	if !exists {
		err = ErrInvalidFilepath.Err(err, "filepath", fileName)
		goto end
	}
end:
	return err
}
