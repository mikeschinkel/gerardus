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
}

func New() *App {
	return &App{}
}

func NewWithDeps(a App) *App {
	return &App{
		dataStore: a.dataStore,
	}
}

func (a *App) DataStore() DataStore {
func (a *App) DataStore() persister.DataStore {
	return a.dataStore
}

func (a *App) Queries() persister.DataStoreQueries {
	return a.dataStore.Queries()
}

func Initialize() {
	Root = New()
	Check.App = Root
}

func DefaultContext() Context {
	return fi.WrapContextFI(context.Background(), FI{
		Persister: PersisterFI{
			InitializeFunc: persister.Initialize,
		},
		Logger: LoggerFI{
			InitializeFunc: logger.Initialize,
		},
	})
}

func (a *App) Main(ctx Context, osArgs []string) (help cli.Help, err error) {
	var invoker *cli.CommandInvoker

	injector := fi.GetFI[FI](ctx)

	err = injector.Logger.Initialize(logger.Params{
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

	a.dataStore, err = injector.Persister.Initialize(ctx,
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
