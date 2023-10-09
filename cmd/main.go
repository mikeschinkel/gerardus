package main

import (
	"context"

	"gerardus/cli"
	"gerardus/collector"
	"gerardus/options"
	"gerardus/persister"
)

func main() {
	err := cli.Initialize()
	if err != nil {
		usage("Failed to initialize; %s", err.Error())
	}
	err = persister.Initialize(context.Background(),
		options.DataFile(),
		collector.SymbolTypes,
	)
	if err != nil {
		usage("Failed to initialize data store; %s", err.Error())
	}
	err = cli.ExecInvokedCommand()
	if err != nil {
		usage("Failed to scan source; %s", err.Error())
	}
}
