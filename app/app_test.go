package app_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/mikeschinkel/gerardus/app"
	"github.com/mikeschinkel/gerardus/cli"
	"github.com/mikeschinkel/gerardus/persister"
	. "github.com/mikeschinkel/go-lib"
	"github.com/mikeschinkel/go-serr"
)

func TestAppMain(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		errStr  string
		stdErr  string
		ctxFunc func(context.Context) context.Context
	}{
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
			stdErr: addProjectGoLangOutput(),
			errStr: "argument cannot be empty [arg_name='<repo_url>']",
		},
		{
			name:   "add project golang https://not.there",
			args:   []string{"add", "project", "golang", "https://not.there"},
			stdErr: addProjectGoLangNotThereOutput(),
			errStr: "not a valid GitHub repo URL [repo_url='https://not.there']",
		},
		//{
		//	name:   "add project golang https://github.com/not/there",
		//	args:   []string{"add", "project", "golang", "https://github.com/not/there"},
		//	stdErr: addProjectGoLangGitHubNotThereOutput(),
		//	errStr: "",
		//	ctxFunc: func(ctx context.Context) context.Context {
		//		return fi.WrapContext(ctx, &app.FI{})
		//	},
		//},
		//{
		//	name:   "add project golang http://github.com/golang/go",
		//	args:   []string{"add", "project", "golang", "http://github.com/golang/go"},
		//	stdErr: addProjectGoLangGitHubGolangGo(),
		//	errStr: "",
		//	ctxFunc: func() context.Context {
		//		return context.WithValue(context.Background(), app.Key, &app.DI{
		//			RepoInfoRequesterFunc: RepoInfoRequesterMock,
		//			UpsertProjectFunc:     UpsertProjectMock,
		//		})
		//	},
		//},
		//{
		//	name:    "add",
		//	args:    []string{"codebase", "golang", "1.21.4"},
		//	errStr: "*",
		//},
		//{
		//	name:    "map",
		//	args:    []string{"map", "golang", "1.21.4"},
		//	errStr: "*",
		//},
	}
	for _, tt := range tests {
		tt.args = RightShift(tt.args, cli.ExecutableFilepath(app.AppName))
		t.Run(tt.name, func(t *testing.T) {
			ctx := app.DefaultContext()
			if tt.ctxFunc != nil {
				ctx = tt.ctxFunc(ctx)
			}
			app.Initialize()
			help, err := app.Root.Main(ctx, tt.args)
			buf := bytes.Buffer{}
			help.Usage(err, &buf)
			if tt.stdErr != buf.String() {
				t.Errorf("Main() value -want +got: %s", cmp.Diff(buf.String(), tt.stdErr))
			}
			if tt.errStr == "" {
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
func CheckURLMock(url string) (err error) {
	switch url {
	case "https://github.com/not/there":
		err = serr.New("oops")
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
func UpsertProjectMock(ctx context.Context, params persister.UpsertProjectParams) (p persister.Project, err error) {
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
func addProjectGoLangOutput() string {
	return `
ERROR: Argument cannot be empty [arg_name='<repo_url>']:
` + projectUsage()
}

func addProjectGoLangNotThereOutput() string {
	return `
ERROR: Not a valid GitHub repo URL [repo_url='https://not.there']:
` + projectUsage()
}
func addProjectGoLangGitHubNotThereOutput() string {
	return `
ERROR: URL could not be dereferenced [repo_url='https://github.com/not/there']:
` + projectUsage()
}

func addProjectGoLangGitHubGolangGo() string {
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
