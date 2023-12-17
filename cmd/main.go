package main

import (
	"os"

	"github.com/mikeschinkel/gerardus/app"
)

func main() {
	app.Initialize()
	help, err := app.Root.Main(app.DefaultContext(), os.Args)
	if err != nil {
		help.Usage(err, os.Stderr)
		os.Exit(1)
	}
}
