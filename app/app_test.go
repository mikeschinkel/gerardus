package app_test

import (
	"bytes"
	"context"
	"database/sql"
	"log/slog"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/mikeschinkel/gerardus/app"
	"github.com/mikeschinkel/gerardus/cli"
	"github.com/mikeschinkel/gerardus/fi"
	"github.com/mikeschinkel/gerardus/logger"
	"github.com/mikeschinkel/gerardus/persister"
	"github.com/mikeschinkel/go-lib"
	"github.com/mikeschinkel/go-serr"
)

type Context = context.Context

type TestOps struct {
	LiveDB bool
}

func TestingContext(tt test, opts TestOps) Context {
	ctx := app.DefaultContext()
	injector := fi.GetFI[app.FI](ctx)
	//--------------
	injector.Logger.InitializeFunc = loggerInitializeMock
	if !opts.LiveDB {
		injector.Persister.InitializeFunc = func(c app.Context, s string, a ...any) (persister.DataStore, error) {
			ds, err := persisterInitializeMock(c, s, a...)
			ds.SetQueries(tt.queries)
			return ds, err
		}
	}
	//--------------
	if tt.fiFunc != nil {
		injector = tt.fiFunc(injector)
	}
	return fi.WrapContextFI(ctx, injector)
}

type test struct {
	name    string
	args    []string
	errStr  string
	stdErr  string
	fiFunc  func(app.FI) app.FI
	queries persister.DataStoreQueries
}

func TestAppMain(t *testing.T) {
	testOpts := TestOps{
		LiveDB: false,
	}
	tests := []test{
		{
			name:   "No CLI arguments",
			args:   []string{},
			stdErr: noCLIArgsOutput(),
			errStr: "no command specified",
		},
		{
			name:   "add",
			args:   []string{"add"},
			stdErr: addArgsOutput(),
			errStr: "no exec func found",
		},
		{
			name:   "add project",
			args:   []string{"add", "project"},
			stdErr: addProjectOutput(),
			errStr: "argument cannot be empty [arg_name='<project>']",
		},
		{
			name:   "add project golang",
			args:   []string{"add", "project", "golang"},
			stdErr: addProjectGolangOutput(),
			errStr: "argument cannot be empty [arg_name='<repo_url>']",
		},
		{
			name:    "add project golang https://not.there",
			args:    []string{"add", "project", "golang", "https://not.there"},
			stdErr:  addProjectGolangNotThereOutput(),
			errStr:  "not a valid GitHub repo URL [repo_url='https://not.there']",
			queries: stubQueriesForLoadProjectByNameMissing(),
		},
		{
			name:    "add project golang http://github.com/not/there",
			args:    []string{"add", "project", "golang", "http://github.com/not/there"},
			stdErr:  addProjectGolangHTTP(),
			errStr:  "repo URL does not begin with https://github.com [repo_url='http://github.com/not/there']",
			queries: stubQueriesForLoadProjectByNameMissing(),
			fiFunc: func(fi app.FI) app.FI {
				fi.CheckURLFunc = CheckURLMock
				return fi
			},
		},
		{
			name:    "add project golang https://not/important",
			args:    []string{"add", "project", "golang", "https://not/important"},
			stdErr:  addProjectGolangAlreadyExists(),
			errStr:  "project found [project='golang']",
			queries: stubQueriesForLoadProjectByName(),
			fiFunc: func(fi app.FI) app.FI {
				fi.CheckURLFunc = CheckURLMock
				return fi
			},
		},
		{
			name:    "add project golang https://github.com/not/there",
			args:    []string{"add", "project", "golang", "https://github.com/not/there"},
			stdErr:  addProjectGolangGitHubNotThereOutput(),
			errStr:  "URL could not be dereferenced [repo_url='https://github.com/not/there']",
			queries: stubQueriesForLoadProjectByNameMissing(),
			fiFunc: func(fi app.FI) app.FI {
				fi.CheckURLFunc = CheckURLMock
				return fi
			},
		},
		//{
		//	name:   "add project golang http://github.com/golang/go",
		//	args:   []string{"add", "project", "golang", "http://github.com/golang/go"},
		//	stdErr: addProjectGolangGitHubGolangGo(),
		//	errStr: "<n/a>",
		//	fiFunc: func(fi app.FI) app.FI {
		//		fi.CheckURLFunc = CheckURLMock
		//		return fi
		//	},
		//},
		//{
		//	name:   "add",
		//	args:   []string{"codebase", "golang", "1.21.4"},
		//	errStr: "*",
		//},
		//{
		//	name:   "map",
		//	args:   []string{"map", "golang", "1.21.4"},
		//	errStr: "*",
		//},
	}
	for _, tt := range tests {
		tt.args = lib.RightShift(tt.args, cli.ExecutableFilepath(app.AppName))
		t.Run(tt.name, func(t *testing.T) {
			ctx := TestingContext(tt, testOpts)
			app.Initialize(ctx)
			root := app.Root
			help, err := root.Main(ctx, tt.args)
			buf := bytes.Buffer{}
			help.Usage(err, &buf)
			if tt.stdErr != buf.String() {
				t.Errorf("Main() value -want +got: %s", cmp.Diff(buf.String(), tt.stdErr))
			}
			if tt.errStr == "<n/a>" {
				return
			}
			if err == nil {
				t.Errorf("Main() error wanting but got no error: %s", tt.errStr)
				return
			}
			if tt.errStr != err.Error() {
				t.Errorf("Main() error -want +got: %s", cmp.Diff(tt.errStr, err.Error()))
			}
		})
	}
}

func stubQueriesForLoadProjectByNameMissing() persister.DataStoreQueries {
	return &app.DataStoreQueriesMock{
		LoadProjectByNameFunc: func(ctx context.Context, name string) (persister.Project, error) {
			return persister.Project{}, sql.ErrNoRows
		},
	}
}

func stubQueriesForLoadProjectByName() persister.DataStoreQueries {
	return &app.DataStoreQueriesMock{
		LoadProjectByNameFunc: func(ctx context.Context, name string) (persister.Project, error) {
			return persister.Project{
				ID:      1,
				Name:    "golang",
				RepoUrl: "https://github.com/golang/go",
			}, nil
		},
	}
}

var sLogger *slog.Logger

func getLogContent() string {
	return sLogger.Handler().(*lib.SLogBufferHandler).Content()
}

func persisterInitializeMock(ctx Context, fp string, types ...any) (ds persister.DataStore, err error) {
	ds = NewTestingDataStore()
	err = ds.Initialize(ctx)
	return ds, err
}

func loggerInitializeMock(logger.Params) error {
	sLogger = slog.New(lib.NewSLogBufferHandler())
	slog.SetDefault(sLogger)
	return nil
}

func CheckURLMock(url string) (err error) {
	switch url {
	case "https://github.com/not/there":
		err = serr.New("oops")
	case "https://github.com/golang/go":
		err = nil
	}
	return err
}

func RepoInfoRequesterMock(url string) (ri persister.RepoInfo, err error) {
	switch url {
	case "https://github.com/not/there":
		err = serr.New("oops")
	case "https://github.com/golang/go":
		ri = persister.RepoInfo{
			Description: "The Go programming language",
			Homepage:    "https://go.dev",
		}
	}
	return ri, err
}

func UpsertProjectMock(ctx Context, params persister.UpsertProjectParams) (p persister.Project, err error) {
	return p, err
}

func noCLIArgsOutput() string {
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

func addArgsOutput() string {
	return `
ERROR: There is no 'add' command, but there are these commands:

  Usage: gerardus [<options>] <command> [<args>]

  Commands:

    - add codebase <project> <version_tag>
    - add project <project> <repo_url> [<about> [<website>]]

    Global Options:

      -data=<data_file>: Data file (sqlite3)
`
}

func addProjectOutput() string {
	return `
ERROR: Argument cannot be empty [arg_name='<project>']:
` + projectUsage()
}
func addProjectGolangOutput() string {
	return `
ERROR: Argument cannot be empty [arg_name='<repo_url>']:
` + projectUsage()
}

func addProjectGolangNotThereOutput() string {
	return `
ERROR: Not a valid GitHub repo URL [repo_url='https://not.there']:
` + projectUsage()
}
func addProjectGolangGitHubNotThereOutput() string {
	return `
ERROR: URL could not be dereferenced [repo_url='https://github.com/not/there']:
` + projectUsage()
}
func addProjectGolangAlreadyExists() string {
	return `
ERROR: Project found [project='golang']:
` + projectUsage()
}
func addProjectGolangHTTP() string {
	return `
ERROR: Repo URL does not begin with https://github.com [repo_url='http://github.com/not/there']:
` + projectUsage()
}

func addProjectGolangGitHubGolangGo() string {
	return ``
}

func projectUsage() string {
	return `
  Usage: gerardus [<options>] <command> [<args>]

  Command:

    - add project <project> <repo_url> [<about> [<website>]]

        Args:

          project:   Project name, e.g. 'golang'
          repo_url:  The full GitHub repository URL for the project, e.g. https://github.com/golang/go
          about:     Repo description. Defaults to 'about' from the GitHub API
          website:   Project website URL. Defaults to 'website' from the GitHub API

        Global Options:

          -data=<data_file>: Data file (sqlite3)
`
}

var _ persister.DataStore = (*TestingDataStore)(nil)

type TestingDataStore struct {
	persister.DataStore
}

func NewTestingDataStore() *TestingDataStore {
	ds := persister.NewSqliteDataStore("/tmp/test.db")
	return &TestingDataStore{
		DataStore: ds,
	}
}

func (db *TestingDataStore) Open() (err error) {
	return nil
}
func (db *TestingDataStore) Queries() persister.DataStoreQueries {
	if db.DataStore.Queries() == nil {
		panic("DatastoreQueries NOT SET for TESTING.")
	}
	return db.DataStore.Queries()
}

//func (t TestingDataStore) Open() error {
//	return nil
//}
//
//func (t TestingDataStore) Query(ctx context.Context, sql string) error {
//	return nil
//}
//
//func (t TestingDataStore) DB() *sql.DB {
//	return &sql.DB{}
//}
//
//func (t TestingDataStore) Initialize(ctx context.Context) error {
//	return nil
//}
