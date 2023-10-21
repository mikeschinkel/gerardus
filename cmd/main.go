package main

import (
	"context"

	"gerardus/cli"
	"gerardus/collector"
	"gerardus/logger"
	"gerardus/options"
	"gerardus/persister"
)

const AppName = "gerardus"

func main() {
	var err error

	err = logger.Initialize(AppName)
	if err != nil {
		usage("%s.", err.Error())
	}
	err = cli.Initialize(AppName)
	if err != nil {
		usage("%s.", err.Error())
	}
	err = persister.Initialize(context.Background(),
		options.DataFile(),
		collector.SymbolTypes,
	)
	if err != nil {
		usage("Failed to initialize data store; %s", err.Error())
	}
	err = cli.ValidateCLIArgs()
	if err != nil {
		usage("%s.", err.Error())
	}
	err = cli.ExecInvokedCommand()
	if err != nil {
		usage("Failed to scan source; %s", err.Error())
	}
}
