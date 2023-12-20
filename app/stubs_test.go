package app_test

import (
	"context"
	"database/sql"
	"log/slog"
	"testing"

	"github.com/mikeschinkel/gerardus/app"
	"github.com/mikeschinkel/gerardus/fi"
	"github.com/mikeschinkel/gerardus/logger"
	"github.com/mikeschinkel/gerardus/persister"
	"github.com/mikeschinkel/go-lib"
)

// UseStubs allows the developer to easily disable stubs for when developing
// tests to witness behavior that needs to be stubbed out. Normally this should
// be set to 'true'; if it has been checked into source code with 'false' that
// would be a mistake.
const UseStubs = true

type Context = context.Context

// TestStub just here to stop Go and Goland from bickering
func TestStub(t *testing.T) {}

type TestOps struct {
	NoStub bool
}

func TestingContext(tt test, opts TestOps) Context {
	ctx := app.DefaultContext()
	if !opts.NoStub {
		injector := fi.GetFI[app.FI](ctx)
		injector.Logger.InitializeFunc = loggerInitializeStub
		injector.Persister.InitializeFunc = func(c app.Context, s string, a ...any) (persister.DataStore, error) {
			ds, err := persisterInitializeStub(c, s, a...)
			ds.SetQueries(tt.queries)
			return ds, err
		}
		if tt.fi != nil {
			injector = tt.fi(injector)
		}
		ctx = fi.WrapContextFI[app.FI](ctx, injector)
	}
	return ctx
}

func SuccessfulUpsertProjectStub(ctx context.Context, arg persister.UpsertProjectParams) (persister.Project, error) {
	return projectID1NameGoLangStub(), nil
}
func LoadFoundProjectByNameStub(ctx context.Context, name string) (persister.Project, error) {
	return projectID1NameGoLangStub(), nil
}
func LoadMissingProjectByNameStub(ctx context.Context, name string) (persister.Project, error) {
	return persister.Project{}, sql.ErrNoRows
}
func projectID1NameGoLangStub() persister.Project {
	return persister.Project{
		ID:      1,
		Name:    "golang",
		RepoUrl: "https://github.com/golang/go",
	}
}

func persisterInitializeStub(ctx Context, fp string, types ...any) (ds persister.DataStore, err error) {
	ds = NewDataStoreStub()
	err = ds.Initialize(ctx)
	return ds, err
}

func stubbedLogContent() string {
	return loggerStub.Handler().(*lib.SLogBufferHandler).Content()
}

var loggerStub *slog.Logger

func loggerInitializeStub(logger.Params) error {
	loggerStub = slog.New(lib.NewSLogBufferHandler())
	slog.SetDefault(loggerStub)
	return nil
}

func CheckURLStub(url string) (err error) {
	switch url {
	case "https://github.com/not/there":
		err = app.ErrURLCouldNotBeDereferenced.Args("repo_url", url)
	case "https://github.com/golang/go":
		err = nil
	}
	return err
}

func RequestGitHubRepoInfoStub(url string) (ri *persister.RepoInfo, err error) {
	switch url {
	case "https://github.com/not/there":
		err = app.ErrURLCouldNotBeDereferenced.Args("repo_url", url)
	case "https://github.com/golang/go":
		ri = &persister.RepoInfo{
			Description: "The Go programming language",
			Homepage:    "https://go.dev",
		}
	}
	return ri, err
}

func UpsertProjectStub(ctx Context, params persister.UpsertProjectParams) (p persister.Project, err error) {
	return p, err
}

var _ persister.DataStore = (*DataStoreStub)(nil)

type DataStoreStub struct {
	persister.DataStore
}

func NewDataStoreStub() *DataStoreStub {
	ds := persister.NewSqliteDataStore("/tmp/test.db")
	return &DataStoreStub{
		DataStore: ds,
	}
}

func (db *DataStoreStub) Open() (err error) {
	return nil
}
func (db *DataStoreStub) Queries() (q persister.DataStoreQueries) {
	q = db.DataStore.Queries()
	if q == nil {
		panic("DatastoreQueries NOT SET for TESTING.")
	}
	return q
}

//goland:noinspection GoUnusedParameter
func (db *DataStoreStub) Query(ctx context.Context, sql string) error {
	return nil
}

func (db *DataStoreStub) DB() *sql.DB {
	return &sql.DB{}
}

//goland:noinspection GoUnusedParameter
func (db *DataStoreStub) Initialize(ctx context.Context) error {
	return nil
}
