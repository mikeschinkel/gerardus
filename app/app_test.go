package app_test

import (
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/mikeschinkel/gerardus/app"
	"github.com/mikeschinkel/gerardus/cli"
	"github.com/mikeschinkel/gerardus/persister"
	"github.com/mikeschinkel/go-lib"
)

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
	//goland:noinspection GoBoolExpressions
	testOpts := TestOps{
		NoStub: !UseStubs,
	}
	tests := []test{
		{
			name:   "FAIL — NO COMMAND",
			fail:   true,
			args:   []string{},
			output: noCLIArgsOutput(),
			errStr: "no command specified",
		},
		{
			name:   "FAIL — NO EXEC FUNC",
			fail:   true,
			args:   []string{"add"},
			output: addArgsOutput(),
			errStr: "no exec func found",
		},
		{
			name:   "FAIL — NO PROJECT ARG",
			fail:   true,
			args:   []string{"add", "project"},
			output: "\nERROR: Argument cannot be empty [arg_name='<project>']:\n" + projectUsage(),
			errStr: "argument cannot be empty [arg_name='<project>']",
		},
		{
			name:   "FAIL — NO REPO URL ARG",
			fail:   true,
			args:   []string{"add", "project", "golang"},
			output: "\nERROR: Argument cannot be empty [arg_name='<repo_url>']:\n" + projectUsage(),
			errStr: "argument cannot be empty [arg_name='<repo_url>']",
		},
		{
			name:   "FAIL — NO GITHUB URL",
			fail:   true,
			args:   []string{"add", "project", "golang", "https://not.there"},
			output: "\nERROR: Not a valid GitHub repo URL [repo_url='https://not.there']:\n" + projectUsage(),
			errStr: "not a valid GitHub repo URL [repo_url='https://not.there']",
			queries: &app.DataStoreQueriesStub{
				LoadProjectByNameFunc: LoadMissingProjectByNameStub,
			},
		},
		{
			name:   "FAIL — NO HTTPS",
			fail:   true,
			args:   []string{"add", "project", "golang", "http://github.com/not/there"},
			output: "\nERROR: Repo URL does not begin with https://github.com [repo_url='http://github.com/not/there']:\n" + projectUsage(),
			errStr: "repo URL does not begin with https://github.com [repo_url='http://github.com/not/there']",
			queries: &app.DataStoreQueriesStub{
				LoadProjectByNameFunc: LoadMissingProjectByNameStub,
			},
		},
		{
			name:   "FAIL — PROJECT EXISTS",
			fail:   true,
			args:   []string{"add", "project", "golang", "https://github.com/golang/go"},
			output: "\nERROR: Project exists [project='golang']:\n" + projectUsage(),
			errStr: "project exists [project='golang']",
			queries: &app.DataStoreQueriesStub{
				LoadProjectByNameFunc: LoadFoundProjectByNameStub,
			},
		},
		{
			name:   "FAIL — CANNOT DEREFERENCE PROJECT URL",
			fail:   true,
			args:   []string{"add", "project", "golang", "https://github.com/not/there"},
			output: "\nERROR: URL could not be dereferenced [repo_url='https://github.com/not/there']:\n" + projectUsage(),
			errStr: "URL could not be dereferenced [repo_url='https://github.com/not/there']",
			queries: &app.DataStoreQueriesStub{
				LoadProjectByNameFunc: LoadMissingProjectByNameStub,
			},
			fi: func(fi app.FI) app.FI {
				fi.Persister.RequestGitHubRepoInfoFunc = RequestGitHubRepoInfoStub // TODO: Verify this is called by this test
				return fi
			},
		},
		{
			name:   "SUCCESS — ADD PROJECT",
			fail:   false,
			args:   []string{"add", "project", "golang", "https://github.com/golang/go"},
			output: "\nSuccessfully added project 'golang' with repo URL https://github.com/golang/go.\n",
			errStr: "<n/a>",
			queries: &app.DataStoreQueriesStub{
				LoadProjectByNameFunc: LoadMissingProjectByNameStub,
				UpsertProjectFunc:     SuccessfulUpsertProjectStub,
			},
		},
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
			buf := bytes.Buffer{}
			cli.StdoutWriter = &buf
			help, err := root.Main(ctx, tt.args)
			if err != nil {
				help.Usage(err, &buf)
			}
			if tt.output != buf.String() {
				t.Errorf("Main() value -want +got: %s", cmp.Diff(tt.output, buf.String()))
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
