package main

import (
	"context"
	"os"

	"github.com/mikeschinkel/gerardus/app"
)

func main() {
	help, err := app.Main(context.TODO(), os.Args, app.MainOpts{})
	if err != nil {
		help.Usage(err, os.Stderr)
		os.Exit(1)
	}
}
