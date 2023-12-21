package app_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/mikeschinkel/gerardus/app"
	"github.com/mikeschinkel/gerardus/persister"
)

func TestAddCodebase(t *testing.T) {}

func addCodebaseTests() []test {
	return []test{
		{
			name:   "FAIL — NO PROJECT ARG",
			fail:   true,
			args:   []string{"add", "codebase"},
			output: "\nERROR: Argument cannot be empty [arg_name='<project>']:\n" + addCodebaseUsage(),
			errStr: "argument cannot be empty [arg_name='<project>']",
		},
		{
			name:   "FAIL — NO REPO URL ARG",
			fail:   true,
			args:   []string{"add", "codebase", "golang"},
			output: "\nERROR: Argument cannot be empty [arg_name='<version_tag>']:\n" + addCodebaseUsage(),
			errStr: "argument cannot be empty [arg_name='<version_tag>']",
		},
		{
			name:   "FAIL — INVALID VERSION TAG",
			fail:   true,
			args:   []string{"add", "codebase", "golang", "foo-bar"},
			output: "\nERROR: Version is not valid [version_tag='foo-bar'] [hint='Version must be semver.org compatible']:\n" + addCodebaseUsage(),
			errStr: "version is not valid [version_tag='foo-bar'] [hint='Version must be semver.org compatible']",
		},
		{
			name:   "FAIL — VERSION_TAG EXISTS",
			fail:   true,
			args:   []string{"add", "codebase", "golang", "go1.21.4"},
			output: "\nERROR: Version already exists [project='golang'] [version_tag='go1.21.4']:\n" + addCodebaseUsage(),
			errStr: "version already exists [project='golang'] [version_tag='go1.21.4']",
			queries: &app.DataStoreQueriesStub{
				LoadProjectByNameFunc:                 LoadFoundProjectByNameStub,
				LoadCodebaseIDByProjectAndVersionFunc: LoadFoundCodebaseIDByProjectAndVersionStub,
			},
		},
		{
			name:   "SUCCESS — ADD CODEBASE",
			fail:   false,
			args:   []string{"add", "codebase", "golang", "go1.21.4"},
			output: "\nSuccessfully added codebase for 'golang' version 'go1.21.4'.\n",
			errStr: "<n/a>",
			queries: &app.DataStoreQueriesStub{
				LoadProjectByNameFunc:                 LoadFoundProjectByNameStub,
				LoadCodebaseIDByProjectAndVersionFunc: LoadMissingCodebaseIDByProjectAndVersionStub,
				UpsertCodebaseFunc:                    SuccessfulUpsertCodebaseStub,
			},
		},
	}
}

func addCodebaseUsage() string {
	return `
  Usage: gerardus [<options>] <command> [<args>]

  Command:

    - add codebase <project> <version_tag>

        Args:

          project:     Project name, e.g. 'golang'
          version_tag: Git version tag

        Global Options:

          -data=<data_file>: Data file (sqlite3)
`
}

func SuccessfulUpsertCodebaseStub(ctx context.Context, arg persister.UpsertCodebaseParams) (persister.Codebase, error) {
	return codebaseID1ProjectGolangVersion1214(), nil
}

func LoadFoundCodebaseIDByProjectAndVersionStub(ctx context.Context, arg persister.LoadCodebaseIDByProjectAndVersionParams) (int64, error) {
	return 1, nil
}
func LoadMissingCodebaseIDByProjectAndVersionStub(ctx context.Context, arg persister.LoadCodebaseIDByProjectAndVersionParams) (int64, error) {
	return 0, sql.ErrNoRows
}

func codebaseID1ProjectGolangVersion1214() persister.Codebase {
	return persister.Codebase{
		ID:         1,
		ProjectID:  1,
		VersionTag: "go1.21.4",
		SourceUrl:  "https://github.com/golang/go/tree/go1.21.4/src",
	}
}
