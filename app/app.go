package app

import (
	"context"

	"github.com/mikeschinkel/gerardus/cli"
	"github.com/mikeschinkel/gerardus/collector"
	"github.com/mikeschinkel/gerardus/logger"
	"github.com/mikeschinkel/gerardus/options"
	"github.com/mikeschinkel/gerardus/parser"
	"github.com/mikeschinkel/gerardus/persister"
)

type MainOpts struct {
}

func Main(ctx context.Context, osArgs []string, mo MainOpts) (help cli.Help, err error) {
	var i *cli.CommandInvoker

	err = logger.Initialize(logger.Params{
		Name:      AppName,
		EnvPrefix: EnvPrefix,
	})
	if err != nil {
		goto end
	}
	err = options.Initialize(options.Params{
		EnvPrefix: EnvPrefix,
	})
	if err != nil {
		goto end
	}
	i, err = cli.Initialize(cli.Params{
		AppName: AppName,
		OSArgs:  osArgs,
	})

	help = cli.NewHelp(i)

	if err != nil {
		goto end
	}
	err = persister.Initialize(ctx,
		options.DataFile(),
		collector.SymbolTypes,
		parser.PackageTypes,
	)
	if err != nil {
		err = ErrFailedToInitDataStore.Err(err, "data_file", options.DataFile())
		goto end
	}
	err = i.Validate(ctx)
	if err != nil {
		goto end
	}
	err = i.InvokeCommand(ctx)
	if err != nil {
		goto end
	}
end:
	return help, err
}
