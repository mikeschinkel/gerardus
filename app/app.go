package app

import (
	"context"

	"github.com/mikeschinkel/gerardus/cli"
	"github.com/mikeschinkel/gerardus/collector"
	"github.com/mikeschinkel/gerardus/fi"
	"github.com/mikeschinkel/gerardus/logger"
	"github.com/mikeschinkel/gerardus/options"
	"github.com/mikeschinkel/gerardus/parser"
	"github.com/mikeschinkel/gerardus/persister"
)

var Root *App

type App struct {
	dataStore persister.DataStore
	fi        FI
}

func New() *App {
	return &App{
		fi: DefaultFI(),
	}
}

func NewWithDeps(a App) *App {
	if a.fi.IsValid() {
		a.fi = New().fi
	}
	return &a
}

func (a *App) DataStore() persister.DataStore {
	return a.dataStore
}

func (a *App) Queries() persister.DataStoreQueries {
	return a.dataStore.Queries()
}

func Initialize(ctx Context) {
	Root = New()
	Root.fi = fi.GetFI[FI](ctx)
	Check.App = Root
}

func DefaultFI() FI {
	return FI{
		Persister: PersisterFI{
			InitializeFunc: persister.Initialize,
		},
		Logger: LoggerFI{
			InitializeFunc: logger.Initialize,
		},
	}
}

func DefaultContext() Context {
	return fi.WrapContextFI(context.Background(), DefaultFI())
}

func (a *App) Main(ctx Context, osArgs []string) (help cli.Help, err error) {
	var invoker *cli.CommandInvoker

	err = a.fi.Logger.Initialize(logger.Params{
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
	invoker, err = cli.Initialize(ctx, cli.Params{
		AppName: AppName,
		OSArgs:  osArgs,
	})
	help = cli.NewHelp(invoker)
	if err != nil {
		goto end
	}

	a.dataStore, err = a.fi.Persister.Initialize(ctx,
		options.DataFile(),
		collector.SymbolTypes,
		parser.PackageTypes,
	)
	if err != nil {
		goto end
	}
	err = invoker.Validate(ctx)
	if err != nil {
		goto end
	}
	err = invoker.InvokeCommand(ctx)
	if err != nil {
		goto end
	}
end:
	return help, err
}
