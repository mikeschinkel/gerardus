package app

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/mikeschinkel/gerardus/cli"
	"github.com/mikeschinkel/gerardus/collector"
	"github.com/mikeschinkel/gerardus/fi"
	"github.com/mikeschinkel/gerardus/logger"
	"github.com/mikeschinkel/gerardus/options"
	"github.com/mikeschinkel/gerardus/parser"
	"github.com/mikeschinkel/gerardus/persister"
	"github.com/mikeschinkel/go-serr"
)

var Root *App = New()

type App struct {
	dataStore persister.DataStore
	project   *persister.Project
	repoInfo  *persister.RepoInfo
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
	Root.fi = fi.GetFI[FI](ctx)
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

func (a *App) checkProjectName(ctx Context, project any, arg *cli.Arg) (err error) {
	var p persister.Project
	var injector FI

	projName := project.(string)
	injector = AssignFI(ctx, FI{Persister: PersisterFI{
		LoadProjectByNameFunc: a.Queries().LoadProjectByName,
	}})
	p, err = injector.Persister.LoadProjectByName(ctx, projName)
	if errors.Is(err, sql.ErrNoRows) {
		err = ErrProjectNotFound.Args("project", project)
		goto end
	}
	if err != nil {
		err = ErrProjectNotFound.Err(err, "project", project)
		goto end
	}
	arg.Message = serr.New("project exists").Args(ProjectArg, projName).Error()
	a.project = &p
end:
	return err
}

func (a *App) checkRepoURL(ctx Context, url any, arg *cli.Arg) (err error) {
	var parts []string
	var numParts int
	var injector FI

	repoURL := url.(string)
	if len(repoURL) == 0 {
		err = ErrNoRepoURLSpecified
		goto end
	}
	if len(repoURL) == 1 && repoURL == "." {
		// A dot/period means ignore the repo
		goto end
	}
	parts = strings.Split(strings.TrimRight(repoURL, "/"), "/")
	numParts = len(parts)
	if numParts != 5 {
		err = ErrInvalidGitHubRepoURL
		goto end
	}
	if strings.Join(parts[:3], "/") != "https://github.com" {
		err = ErrInvalidGitHubRepoRootURL
		goto end
	}
	injector = AssignFI(ctx, FI{CheckURLFunc: cli.CheckURL})
	err = injector.CheckURL(repoURL)
	if err != nil {
		err = ErrURLCouldNotBeDereferenced
		goto end
	}
	injector = AssignFI(ctx, FI{Persister: PersisterFI{RepoInfoRequesterFunc: persister.RequestGitHubRepoInfo}})
	a.repoInfo, err = injector.Persister.RepoInfoRequesterFunc(repoURL)
	if err != nil {
		goto end
	}

end:
	if err != nil && len(repoURL) > 0 {
		err = err.(serr.SError).Args("repo_url", repoURL)
	}
	return err
}

func (a *App) checkVersionTag(ctx Context, tag any, arg *cli.Arg) (err error) {
	var verTag string
	var injector FI

	projName := options.ProjectName()
	if len(projName) == 0 {
		err = ErrNoProjectSpecified
		goto end
	}
	verTag = tag.(string)
	if len(verTag) == 0 {
		err = ErrNoVersionTagSpecified.Args("project", projName)
		goto end
	}

	injector = AssignFI(ctx, FI{Persister: PersisterFI{
		LoadCodebaseIDByProjectAndVersionFunc: a.Queries().LoadCodebaseIDByProjectAndVersion,
	}})
	_, err = injector.Persister.LoadCodebaseIDByProjectAndVersion(ctx, persister.LoadCodebaseIDByProjectAndVersionParams{
		Name:       projName,
		VersionTag: verTag,
	})
	if err != nil {
		err = ErrFailedToAddCodebase.Err(err, "project", projName, "version_tag", verTag)
		goto end
	}
end:
	return err
}
