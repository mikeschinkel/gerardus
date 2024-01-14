package test

import (
	"fmt"
	"testing"

	"github.com/mikeschinkel/gerardus/app"
	"github.com/mikeschinkel/go-serr"
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
			output: "\nERROR: Version tag does not exist [version_tag='foo-bar'] [project='golang']:\n" + mapUsage(),
			errStr: "version tag does not exist [version_tag='foo-bar'] [project='golang']",
			queries: &app.DataStoreQueriesStub{
				LoadProjectByNameFunc:                 LoadFoundProjectByNameStub,
				LoadCodebaseIDByProjectAndVersionFunc: LoadMissingCodebaseIDByProjectAndVersionStub,
			},
		},
		{
			name:   "FAIL — MAP GOLANG v1.21.4",
			fail:   true,
			args:   []string{"map", "golang", "v1.21.4"},
			output: "Scanning Go source at " + TestSourceDir + "...\n\nERROR: `map` command failed [source_dir='/tmp/test/dir']:\n" + mapUsage(),
			errStr: "`map` command failed [source_dir='/tmp/test/dir']",
			queries: &app.DataStoreQueriesStub{
				LoadProjectByNameFunc:                 LoadFoundProjectByNameStub,
				LoadCodebaseIDByProjectAndVersionFunc: LoadFoundCodebaseIDByProjectAndVersionStub,
			},
			fi: func(fi app.FI) app.FI {
				fi.App.Map = mapFailureStub
				return fi
			},
		},
		{
			name:   "SUCCESS — MAP GOLANG v1.21.4",
			fail:   false,
			args:   []string{"map", "golang", "v1.21.4"},
			output: "Scanning Go source at " + TestSourceDir + "...\nScanning completed successfully.",
			errStr: "<n/a>",
			queries: &app.DataStoreQueriesStub{
				LoadProjectByNameFunc:                 LoadFoundProjectByNameStub,
				LoadCodebaseIDByProjectAndVersionFunc: LoadFoundCodebaseIDByProjectAndVersionStub,
			},
			fi: func(fi app.FI) app.FI {
				fi.App.Map = mapSuccessStub
				return fi
			},
		},
	}
}

func mapSuccessStub(Context, string, string, *app.App) error {
	return nil
}

func mapFailureStub(Context, string, string, *app.App) error {
	return serr.New("test-invoked failure")
}

func mapUsage() string {
	return fmt.Sprintf(`
  Usage: gerardus [<options>] <command> [<args>]

  Command:

    - map [-src=<source_dir>] <project> <version_tag>

        Options:

          -src=<source_dir>: Source directory
           Default:  %s

        Args:

          project:     Project name, e.g. 'golang'
          version_tag: Git version tag

        Global Options:

          -data=<data_file>: Data file (sqlite3)
`, TestSourceDir)
}
