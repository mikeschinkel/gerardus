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
		{
			name:   "FAIL — INVALID VERSION TAG",
			fail:   true,
			args:   []string{"map", "golang", "foo-bar"},
			output: "\nERROR: Version tag does not exist [version_tag='go1.21.4'] [project='golang']:\n" + mapUsage(),
			errStr: "version tag does not exist [version_tag='go1.21.4'] [project='golang']",
			queries: &app.DataStoreQueriesStub{
				LoadProjectByNameFunc:                 LoadFoundProjectByNameStub,
				LoadCodebaseIDByProjectAndVersionFunc: LoadMissingCodebaseIDByProjectAndVersionStub,
			},
		},
	}
}

func mapUsage() string {
	return fmt.Sprintf(`
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
`, TestSourceDir)
}
