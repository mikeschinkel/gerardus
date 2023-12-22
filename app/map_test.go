package app_test

import (
	"testing"
)

func TestMapCommand(t *testing.T) {}

//goland:noinspection GoUnusedFunction
func mapTests() []test {
	return []test{
		{
			name:   "FAIL — NO PROJECT",
			fail:   true,
			args:   []string{"map"},
			output: mapArgsOutput(),
			errStr: "argument cannot be empty [arg_name='<project>']",
		},
	}
}

func mapArgsOutput() string {
	return `
ERROR: Argument cannot be empty [arg_name='<project>']:

  Usage: gerardus [<options>] <command> [<args>]

  Command:

    - map [-src=<source_dir>] <project> <version_tag>

        Options:

          -src=<source_dir>: Source directory
           Default:  /Users/mikeschinkel/Projects/gerardus

        Args:

          project:     Project name, e.g. 'golang'
          version_tag: Git version tag

        Global Options:

          -data=<data_file>: Data file (sqlite3)
`
}
