package test

import (
	"testing"
)

func TestAddCommand(t *testing.T) {}

//goland:noinspection GoUnusedFunction
func addTests() []test {
	return []test{
		{
			name:   "FAIL — NO EXEC FUNC",
			fail:   true,
			args:   []string{"add"},
			output: addArgsOutput(),
			errStr: "no exec func found",
		},
	}
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
