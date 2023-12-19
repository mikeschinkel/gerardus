package main

import (
	"os"

	"github.com/mikeschinkel/gerardus/app"
)

func main() {
	ctx := app.DefaultContext()
	app.Initialize(ctx)
	help, err := app.Root.Main(ctx, os.Args)
	if err != nil {
		help.Usage(err, os.Stderr)
		os.Exit(1)
	}
}
