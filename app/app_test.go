package app_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/mikeschinkel/gerardus/app"
	"github.com/mikeschinkel/gerardus/cli"
	. "github.com/mikeschinkel/go-lib"
)

func TestAppMain(t *testing.T) {
	tests := []struct {
		name   string
		args   []string
		errStr string
		stdErr string
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
		//{
		//	name:    "add project golang",
		//	args:    []string{"add", "project", "golang"},
		//	errStr: "*",
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
			help, err := app.Main(context.Background(), tt.args, app.MainOpts{})
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

  Usage: gerardus [<options>] <command> [<args>]

  Command:

    - add project <project> <repo_url> [<about> [<website>]]

        Args:

          project:   Project name, e.g. 'golang'
          repo_url:  The full GitHub repository URL for the project, e.g. https://github.com/golang/go
          about:
          website:

        Global Options:

          -data=<data_file>: Data file (sqlite3)
`
}
