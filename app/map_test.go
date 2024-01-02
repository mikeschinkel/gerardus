package app_test

import (
	"testing"

	"github.com/mikeschinkel/gerardus/app"
)

func TestMapCommand(t *testing.T) {}

//goland:noinspection GoUnusedFunction
func mapTests() []test {
	return []test{
		{
			name:   "FAIL — NO PROJECT",
			fail:   true,
			args:   []string{"map"},
			output: "\nERROR: Value cannot be empty [arg_name=project]:\n" + mapUsage(),
			errStr: "value cannot be empty [arg_name=project]",
		},
		{
			name:   "FAIL — INVALID PROJECT",
			fail:   true,
			args:   []string{"map", "foobar", "v1.2.3"},
			output: "\nERROR: Project not found [project='foobar']:\n" + mapUsage(),
			errStr: "project not found [project='foobar']",
			queries: &app.DataStoreQueriesStub{
				LoadProjectByNameFunc: LoadMissingProjectByNameStub,
			},
		},
	}
}

func mapUsage() string {
	return `
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
