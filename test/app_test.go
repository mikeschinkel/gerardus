package test

import (
	"bytes"
	"context"
	"database/sql"
	"log/slog"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/mikeschinkel/gerardus/app"
	"github.com/mikeschinkel/gerardus/cli"
	"github.com/mikeschinkel/gerardus/fi"
	"github.com/mikeschinkel/gerardus/logger"
	"github.com/mikeschinkel/gerardus/options"
	"github.com/mikeschinkel/gerardus/persister"
	"github.com/mikeschinkel/go-diffator"
	"github.com/mikeschinkel/go-lib"
)

// UseStubs allows the developer to easily disable stubs for when developing
// tests to witness behavior that needs to be stubbed out. Normally this should
// be set to 'true'; if it has been checked into source code with 'false' that
// would be a mistake.
const (
	UseStubs      = true
	TestSourceDir = "/tmp/test/dir"
)

type Context = context.Context

type TestOps struct {
	NoStub bool
}

type test struct {
	name    string
	args    []string
	output  string
	errStr  string
	fail    bool
	fi      func(app.FI) app.FI
	queries persister.DataStoreQueries
}

func TestAppMain(t *testing.T) {
	t.Run("CLI Tests", func(t *testing.T) {
		t.Run("Root Tests", func(t *testing.T) {
			runTests(t, rootTests())
		})
		t.Run("Add Tests", func(t *testing.T) {
			runTests(t, addTests())
		})
		t.Run("Add Project Tests", func(t *testing.T) {
			runTests(t, addProjectTests())
		})
		t.Run("Add Codebase Tests", func(t *testing.T) {
			runTests(t, addCodebaseTests())
		})
		t.Run("Map Tests", func(t *testing.T) {
			runTests(t, mapTests())
		})
		t.Run("Help Tests", func(t *testing.T) {
			runTests(t, helpTests())
		})
	})
}

func runTests(t *testing.T, tests []test) {
	//goland:noinspection GoBoolExpressions
	testOpts := TestOps{
		NoStub: !UseStubs,
	}
	opts := &diffator.StringOpts{
		MinSubstrLen: diffator.Int(2),
	}
	for _, tt := range tests {
		tt.args = lib.RightShift(tt.args, cli.ExecutableFilepath(app.AppName))
		t.Run(tt.name, func(t *testing.T) {
			ctx := ContextStub(tt, testOpts)
			app.Initialize(ctx)
			root := app.Root
			buf := bytes.Buffer{}
			cli.StdoutWriter = &buf
			cli.StderrWriter = &buf
			help, err := root.Main(ctx, tt.args)
			if err != nil {
				help.Usage(err, &buf)
			}
			if tt.output != buf.String() {
				t.Errorf("Main() value -want +got: %s", cmp.Diff(tt.output, buf.String()))
			}
			if tt.errStr == "<n/a>" && !tt.fail {
				return
			}
			if err == nil && tt.fail {
				t.Errorf("Main() error wanting but got no error: %s", tt.errStr)
				return
			}
			if err != nil && tt.errStr != err.Error() {
				t.Errorf("Main() error: got/want: %s",
					diffator.CompareStrings(tt.errStr, err.Error(), opts),
				)
			}
		})
	}
}

func rootTests() []test {
	return []test{
		{
			name:   "FAIL — NO COMMAND",
			fail:   true,
			args:   []string{},
			output: rootArgsOutput(),
			errStr: "no command specified",
		},
		{
			name:   "FAIL — NO DB FILE",
			fail:   true,
			args:   []string{"-data"},
			output: rootArgsOutput(),
			errStr: "no command specified",
		},
		{
			name:   "FAIL — MISSING DB FILE",
			fail:   true,
			args:   []string{"-data=/not/there.db"},
			output: rootArgsOutput(),
			errStr: "no command specified",
		},
	}
}

func rootArgsOutput() string {
	return `
ERROR: No command specified:

  Usage: gerardus [<options>] <command> [<args>]

  Commands:

    - add codebase <project> <version_tag>
    - add project <project> <repo_url> [<about> [<website>]]
    - help [<command>]
    - map [-src=<source_dir>] <project> <version_tag>

    Global Options:

      -data=<data_file>: Data file (sqlite3)
`
}

func ContextStub(tt test, opts TestOps) Context {
	ctx := app.DefaultContext()
	if !opts.NoStub {
		injector := fi.GetFI[app.FI](ctx)
		injector.App.Map = app.Map
		injector.Logger.Initialize = loggerInitializeStub
		injector.Persister.Initialize = func(c app.Context, s string, a ...any) (persister.DataStore, error) {
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

//goland:noinspection GoUnusedParameter
func persisterInitializeStub(ctx Context, fp string, types ...any) (ds persister.DataStore, err error) {
	// This is just a convenient place to set the `source_dir`
	dir := cli.NewValue(reflect.String, TestSourceDir)
	options.SetSourceDir(dir)
	app.CmdMap.SetFlagDefault("source_dir", dir.String())

	ds = NewDataStoreStub()
	err = ds.Initialize(ctx)
	return ds, err
}

//goland:noinspection GoUnusedFunction
func stubbedLogContent() string {
	return loggerStub.Handler().(*lib.SLogBufferHandler).Content()
}

var loggerStub *slog.Logger

func loggerInitializeStub(logger.Params) error {
	loggerStub = slog.New(lib.NewSLogBufferHandler())
	slog.SetDefault(loggerStub)
	return nil
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
		panic("TODO: tt.queries NOT *YET* SET for this test case.")
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
