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
			output: noCLIArgsOutput(),
			errStr: "no command specified",
		},
		{
			name:   "add",
			args:   []string{"add"},
			output: addArgsOutput(),
			errStr: "no exec func found",
		},
		{
			name:   "add project",
			args:   []string{"add", "project"},
			output: "\nERROR: Argument cannot be empty [arg_name='<project>']:\n" + projectUsage(),
			errStr: "argument cannot be empty [arg_name='<project>']",
		},
		{
			name:   "add project golang",
			args:   []string{"add", "project", "golang"},
			output: "\nERROR: Argument cannot be empty [arg_name='<repo_url>']:\n" + projectUsage(),
			errStr: "argument cannot be empty [arg_name='<repo_url>']",
		},
		{
			name:    "add project golang https://not.there",
			args:    []string{"add", "project", "golang", "https://not.there"},
			stdErr:  addProjectGolangNotThereOutput(),
			errStr:  "not a valid GitHub repo URL [repo_url='https://not.there']",
			queries: stubQueriesForLoadProjectByNameMissing(),
			output: "\nERROR: Not a valid GitHub repo URL [repo_url='https://not.there']:\n" + projectUsage(),
		},
		{
			name:    "add project golang http://github.com/not/there",
			args:    []string{"add", "project", "golang", "http://github.com/not/there"},
			stdErr:  addProjectGolangHTTP(),
			errStr:  "repo URL does not begin with https://github.com [repo_url='http://github.com/not/there']",
			queries: stubQueriesForLoadProjectByNameMissing(),
			fiFunc: func(fi app.FI) app.FI {
				fi.CheckURLFunc = CheckURLMock
			output: "\nERROR: Repo URL does not begin with https://github.com [repo_url='http://github.com/not/there']:\n" + projectUsage(),
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
			output: "\nERROR: Project exists [project='golang']:\n" + projectUsage(),
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
			output: "\nERROR: URL could not be dereferenced [repo_url='https://github.com/not/there']:\n" + projectUsage(),
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
